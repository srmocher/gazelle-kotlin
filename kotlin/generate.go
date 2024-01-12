package kotlin

import (
	"fmt"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/language/proto"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/srmocher/gazelle-kotlin/kotlin/internal"
	"github.com/srmocher/gazelle-kotlin/kotlin/protobuf"
	"log"
	"sort"
	"strings"
)

func filterFilesBySuffix(files []string, suffix string) []string {
	var filteredFiles []string
	for _, file := range files {
		if strings.HasSuffix(file, suffix) {
			filteredFiles = append(filteredFiles, file)
		}
	}
	return filteredFiles
}

type packageImportInfo struct {
	PackageName     string
	ExternalImports []*internal.MavenArtifact
	OtherImports    []string
}

func (*kotlinLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	c := args.Config
	kc := GetKotlinConfig(c)

	var protoRules []*rule.Rule
	var protoImportInfos []*packageImportInfo
	var err error
	if kc.ktGenerateProto {
		protoRules, protoImportInfos, err = genProtoRules(args.Rel, args.OtherGen)
		if err != nil {
			log.Print(err)
		}
	}
	if !kc.buildFileGenerationEnabled {
		// Build file generation is disabled
		return language.GenerateResult{}
	}

	ktFiles := filterFilesBySuffix(args.RegularFiles, ".kt")

	rules, packageImportInfos, err := buildKotlinPackage(args.Rel, kc, ktFiles)
	if err != nil {
		log.Print(err)
		return language.GenerateResult{}
	}

	imports := make([]interface{}, 0, len(packageImportInfos)+len(protoImportInfos))
	for _, pii := range packageImportInfos {
		imports = append(imports, pii)
	}

	for _, pii := range protoImportInfos {
		imports = append(imports, pii)
	}

	return language.GenerateResult{
		Gen:     append(rules, protoRules...),
		Imports: imports,
	}
}

func filterImports(sis []*protobuf.SourceFileInfo, mii *internal.MavenInstallInfo) ([]*internal.MavenArtifact, []string) {
	externalImportsMap := make(map[string]*internal.MavenArtifact)
	otherImportsMap := make(map[string]bool)
	for _, si := range sis {
		for _, imp := range si.GetImports() {
			ma := mii.GetMavenArtifactFromImport(imp)
			// not an external dep
			if ma == nil {
				otherImportsMap[imp] = true
			} else {
				externalImportsMap[imp] = ma
			}
		}
	}

	externalImports := make([]*internal.MavenArtifact, 0, len(externalImportsMap))
	otherImports := make([]string, 0, len(otherImportsMap))
	for _, ma := range externalImportsMap {
		externalImports = append(externalImports, ma)
	}

	for imp := range otherImportsMap {
		otherImports = append(otherImports, imp)
	}

	return externalImports, otherImports
}

func getJavaPackage(rel string, protoPkg proto.Package) string {
	var javaPkg string = ""
	for k, v := range protoPkg.Options {
		if k == "java_package" {
			javaPkg = v
		}
	}

	if javaPkg == "" {
		javaPkg = fmt.Sprintf("%s.%s", rel, protoPkg.Name)
	}

	return javaPkg
}

func genProtoRules(rel string, otherGen []*rule.Rule) ([]*rule.Rule, []*packageImportInfo, error) {
	var rules []*rule.Rule
	var pkgImportInfos []*packageImportInfo

	var protoPkg proto.Package
	foundProto := false

	for _, r := range otherGen {
		if r.Kind() == "proto_library" {
			protoPkg = r.PrivateAttr(proto.PackageKey).(proto.Package)
			foundProto = true
			break
		}
	}
	if !foundProto {
		return []*rule.Rule{}, []*packageImportInfo{}, nil
	}

	javaPkg := getJavaPackage(rel, protoPkg)
	if protoPkg.HasServices {
		grpcRule := rule.NewRule("kt_jvm_grpc_library", "kt_grpc_library")
		grpcRule.SetAttr("visibility", []string{"//visibility:public"})
		pkgImportInfos = append(pkgImportInfos, &packageImportInfo{
			PackageName: javaPkg,
		})
		rules = append(rules, grpcRule)
	}

	protoRule := rule.NewRule("kt_jvm_proto_library", "kt_proto_library")
	protoRule.SetAttr("visibility", []string{"//visibility:public"})
	pkgImportInfos = append(pkgImportInfos, &packageImportInfo{
		PackageName: javaPkg,
	})
	rules = append(rules, protoRule)
	return rules, pkgImportInfos, nil
}

func buildKotlinPackage(rel string, kc *KotlinConfig, ktFiles []string) ([]*rule.Rule, []*packageImportInfo, error) {
	sourceFileInfos, err := kc.kotlinParser.ParseKotlinFiles(ktFiles)
	if err != nil {
		return nil, nil, err
	}

	sourcePackage := ""
	for _, si := range sourceFileInfos {
		if sourcePackage != "" && si.GetPackageName() != sourcePackage {
			return nil, nil, fmt.Errorf("source files in %s dir have conflicting package names: (%s, %s)", rel, sourcePackage, si.GetPackageName())
		} else {
			sourcePackage = si.GetPackageName()
		}
	}

	r := rule.NewRule("kt_jvm_library", defaultKtNamingConvention)
	sort.Strings(ktFiles)
	r.SetAttr("srcs", ktFiles)
	r.SetPrivateAttr("package", sourcePackage)

	if len(kc.ktVisibility) > 0 {
		r.SetAttr("visibility", kc.ktVisibility)
	} else {
		r.SetAttr("visibility", []string{"//visibility:public"})
	}

	externalImports, otherImports := filterImports(sourceFileInfos, kc.mavenInstallInfo)
	return []*rule.Rule{r}, []*packageImportInfo{&packageImportInfo{
		PackageName:     sourcePackage,
		ExternalImports: externalImports,
		OtherImports:    otherImports,
	}}, nil
}
