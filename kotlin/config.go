package kotlin

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
