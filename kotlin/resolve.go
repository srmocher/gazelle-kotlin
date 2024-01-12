package kotlin

import (
	"fmt"
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"log"
	"strings"
)

func (*kotlinLang) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	if !isKtJvmLibrary(r.Kind()) || !isKtJvmProtoLibrary(r.Kind()) {
		return nil
	}

	pkg := r.PrivateAttr("package")
	if pkg != nil {
		pkgStr := pkg.(string)
		return []resolve.ImportSpec{
			{
				Lang: kotlinName,
				Imp:  pkgStr,
			},
		}
	} else {
		return []resolve.ImportSpec{}
	}

}

func (*kotlinLang) Embeds(r *rule.Rule, from label.Label) []label.Label {
	return nil
}

func isStandardImport(imp string) bool {
	return strings.HasPrefix(imp, "kotlin.")
}

func (*kotlinLang) Resolve(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imports interface{}, from label.Label) {
	if !isKtJvmLibrary(r.Kind()) || !isKtJvmProtoLibrary(r.Kind()) {
		return
	}

	kc := GetKotlinConfig(c)

	r.DelAttr("deps")
	pkgImportInfo := imports.(packageImportInfo)

	var deps []string
	for _, imp := range pkgImportInfo.ExternalImports {
		deps = append(deps, kc.mavenInstallInfo.CoordToBazelLabel(imp.Coord))
	}

	for _, imp := range pkgImportInfo.OtherImports {
		// No need for an explicit dep if it's a standard lib import
		if isStandardImport(imp) {
			continue
		}

		impParts := strings.Split(imp, ".")
		impPkg := strings.Join(impParts[:len(impParts)-1], ".")
		l, err := resolveWithIndex(ix, c, impPkg)
		if err != nil {
			log.Print(err)
			continue
		}
		deps = append(deps, l.String())
	}
	r.SetAttr("deps", deps)
}

func resolveWithIndex(ix *resolve.RuleIndex, c *config.Config, imp string) (label.Label, error) {
	matches := ix.FindRulesByImportWithConfig(c, resolve.ImportSpec{Lang: kotlinName, Imp: imp}, kotlinName)
	if len(matches) == 0 {
		return label.NoLabel, fmt.Errorf("no label found for %s", imp)
	}

	if len(matches) > 1 {
		return label.NoLabel, fmt.Errorf("more than 1 match for %s", imp)
	}
	return matches[0].Label, nil
}

func isKtJvmLibrary(kind string) bool {
	return kind == "kt_jvm_library"
}

func isKtJvmProtoLibrary(kind string) bool {
	return kind == "kt_jvm_proto_library" || kind == "kt_jvm_grpc_library"
}
