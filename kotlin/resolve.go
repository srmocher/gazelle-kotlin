package kotlin

import (
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func (*kotlinLang) Imports(c *config.Config, r *rule.Rule, f *rule.File) []ImportSpec {
	return nil
}

func (*kotlinLang) Embeds(r *rule.Rule, from label.Label) []label.Label {
	return nil
}

func (*kotlinLang) Resolve(c *config.Config, ix *RuleIndex, rc *repo.RemoteCache, r *rule.Rule, imports interface{}, from label.Label) {
	return
}
