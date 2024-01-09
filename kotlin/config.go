package kotlin

import (
	"flag"
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	defaultRulesKotlinRepoName = "io_bazel_rules_kotlin"
)

type kotlinConfig struct {
	mavenInstallJsonFilePath string

	buildFileGenerationEnabled bool

	ktNamingConvention string

	ktVisibility []string

	ktGenerateProto bool

	mavenInstallRepoName string

	rulesKotlinRepoName string
}

func newKotlinConfig() *kotlinConfig {
	kc := &kotlinConfig{
		rulesKotlinRepoName: defaultRulesKotlinRepoName,
		ktGenerateProto:     false,
	}
}

func (*kotlinLang) KnownDirectives() []string {
	return []string{
		"kt_generate_proto",
		"kt_naming_convention",
		"kt_visibility",
		"kt_maven_install_json_file",
		"kt_enabled",
		"kt_maven_install_repo",
	}
}

func (*kotlinLang) Configure(c *config.Config, rel string, f *rule.File) {
	var kc *kotlinConfig
	if raw, ok := c.Exts[kotlinName]; !ok {
		gc = newKotlinConfig()
	} else {
		gc = raw.(*kotlinConfig).clone()
	}

	c.Exts[kotlinName] = kc
	if kc.rulesKotlinRepoName == "" {
		kc.rulesKotlinRepoName = defaultRulesKotlinRepoName
	}
}

func (*kotlinLang) CheckFlags(fs *flag.FlagSet, c *Config) error {
	return nil
}

func (*kotlinLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *Config) {
	return nil
}
