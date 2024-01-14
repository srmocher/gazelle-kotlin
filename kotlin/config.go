package kotlin

import (
	"flag"
	"fmt"
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/srmocher/gazelle-kotlin/kotlin/internal"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	defaultRulesKotlinRepoName = "io_bazel_rules_kotlin"
	defaultKtNamingConvention  = "kt_default_library"
)

type KotlinConfig struct {
	// if generation of build targets is enabled at all
	buildFileGenerationEnabled bool

	// the naming convention for the kt_jvm* targets
	ktNamingConvention string

	// the visibility of the generated targets
	ktVisibility []string

	// if generation of kt_jvm_proto* or kt_jvm_grpc* targets is enabled
	// for proto targets
	ktGenerateProto bool

	// whether to include .java files into the targets srcs
	javaEnabled bool

	// the repository name of rules_kotlin which is used to load the rules
	rulesKotlinRepoName string

	kotlinParser *internal.KotlinParser

	mavenInstallInfo *internal.MavenInstallInfo
}

func newKotlinConfig() *KotlinConfig {
	return &KotlinConfig{
		rulesKotlinRepoName: defaultRulesKotlinRepoName,
		ktGenerateProto:     false,
		javaEnabled:         false,
		ktNamingConvention:  defaultKtNamingConvention,
	}
}

func (*kotlinLang) KnownDirectives() []string {
	return []string{
		"kt_generate_proto",
		"java_enabled",
		"kt_naming_convention",
		"kt_visibility",
		"kt_maven_repo",
		"kt_enabled",
	}
}

func (kc *KotlinConfig) clone() *KotlinConfig {
	return kc
}

func parseMavenRepoDirective(repoRoot string, mavenRepoDirective string) (*internal.MavenInstallInfo, error) {
	parts := strings.Split(mavenRepoDirective, " ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("kt_maven_repo directive should specify the maven repo name followed by the repo relative path to the maven_install json file")
	}
	return internal.NewMavenInstallInfo(filepath.Join(repoRoot, parts[1]), parts[0]), nil
}

func (*kotlinLang) Configure(c *config.Config, rel string, f *rule.File) {
	var kc *KotlinConfig
	if raw, ok := c.Exts[kotlinName]; !ok {
		kc = newKotlinConfig()
	} else {
		kc = raw.(*KotlinConfig).clone()
	}

	c.Exts[kotlinName] = kc
	if kc.rulesKotlinRepoName == "" {
		kc.rulesKotlinRepoName = defaultRulesKotlinRepoName
	}

	if f != nil {
		for _, directive := range f.Directives {
			switch directive.Key {
			case "kt_generate_proto":
				ktProtoEnabled, err := strconv.ParseBool(directive.Value)
				if err != nil {
					log.Print(err)
					continue
				}
				kc.ktGenerateProto = ktProtoEnabled
			case "java_enabled":
				javaEnabled, err := strconv.ParseBool(directive.Value)
				if err != nil {
					log.Print(err)
					continue
				}
				kc.javaEnabled = javaEnabled
			case "kt_maven_repo":
				mii, err := parseMavenRepoDirective(c.RepoRoot, directive.Value)
				if err != nil {
					log.Print(err)
					continue
				}
				kc.mavenInstallInfo = mii
			case "kt_enabled":
				ktEnabled, err := strconv.ParseBool(directive.Value)
				if err != nil {
					log.Print(err)
					continue
				}
				kc.buildFileGenerationEnabled = ktEnabled
			case "kt_visibility":
				ktVisibility := strings.Split(directive.Value, ",")
				kc.ktVisibility = ktVisibility
			}
		}
	}
}

func GetKotlinConfig(c *config.Config) *KotlinConfig {
	kc := c.Exts[kotlinName]
	if kc == nil {
		return nil
	}
	return kc.(*KotlinConfig)
}

func (kl *kotlinLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	var kc *KotlinConfig
	var kp *internal.KotlinParser
	var err error
	if raw, ok := c.Exts[kotlinName]; !ok {
		kc = newKotlinConfig()
		c.Exts[kotlinName] = kc
	} else {
		kc = raw.(*KotlinConfig).clone()
	}

	if kc.kotlinParser == nil {
		kp, err = internal.NewKotlinParser()
		if err != nil {
			return err
		}
	}
	kc.kotlinParser = kp
	kl.kotlinParser = kp
	return nil
}

func (*kotlinLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	kc := newKotlinConfig()
	c.Exts[kotlinName] = kc
}
