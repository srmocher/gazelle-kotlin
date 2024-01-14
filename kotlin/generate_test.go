package kotlin

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/language/proto"
	"github.com/bazelbuild/bazel-gazelle/merger"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/bazel-gazelle/testtools"
	"github.com/bazelbuild/bazel-gazelle/walk"
	bzl "github.com/bazelbuild/buildtools/build"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testConfig(t *testing.T, repoRoot string) (*config.Config, []language.Language, []config.Configurer) {
	cexts := []config.Configurer{
		&config.CommonConfigurer{},
		&walk.Configurer{},
		&resolve.Configurer{},
	}
	lang := NewLanguage()
	pl := proto.NewLanguage()
	langs := []language.Language{pl, lang}
	c := testtools.NewTestConfig(t, cexts, langs, []string{
		"-build_file_name=BUILD.old",
		"-repo_root=" + repoRoot,
	})
	cexts = append(cexts, pl, lang)
	return c, langs, cexts
}

func isTest(regularFiles []string) bool {
	for _, name := range regularFiles {
		if name == "BUILD.want" {
			return true
		}
	}
	return false
}

func TestGenerateRules(t *testing.T) {
	c, langs, cexts := testConfig(t, "testdata")
	content := []byte(`
# gazelle:follow **
# gazelle:kt_enabled true
# gazelle:kt_generate_proto true
# gazelle:kt_maven_repo maven_install.json maven
`)
	f, err := rule.LoadData(filepath.FromSlash("BUILD.config"), "config", content)
	if err != nil {
		t.Fatal(err)
	}
	for _, cext := range cexts {
		cext.Configure(c, "", f)
	}
	type testCase struct {
		Name         string
		RulesCount   map[string]int
		ErrorMessage string
	}
	testcases := map[string]testCase{
		"simple/src/main/com/foo": {
			Name: "Simple with no protos",
			RulesCount: map[string]int{
				"kt_jvm_library": 1,
			},
			ErrorMessage: "failed test case for simple package",
		},
		"lib_with_protos/src/main/com/foo": {
			Name: "Package with protos",
			RulesCount: map[string]int{
				"kt_jvm_library":       1,
				"proto_library":        1,
				"kt_jvm_proto_library": 1,
			},
			ErrorMessage: "failed test case for package with protos",
		},
	}

	finishable := langs[1].(language.FinishableLanguage)
	defer finishable.DoneGeneratingRules()
	var configurers []config.Configurer
	for _, lang := range langs {
		configurers = append(configurers, lang)
	}
	walk.Walk(c, configurers, []string{"testdata"}, walk.VisitAllUpdateSubdirsMode, func(dir, rel string, c *config.Config, update bool, oldFile *rule.File, subdirs, regularFiles, genFiles []string) {
		if !isTest(regularFiles) {
			return
		}
		t.Run(rel, func(t *testing.T) {
			var gen []*rule.Rule
			var empty []*rule.Rule
			var loads []rule.LoadInfo
			for _, lang := range langs {
				res := lang.GenerateRules(language.GenerateArgs{
					Config:       c,
					Dir:          dir,
					Rel:          rel,
					File:         oldFile,
					Subdirs:      subdirs,
					RegularFiles: regularFiles,
					GenFiles:     genFiles,
					OtherGen:     gen,
					OtherEmpty:   empty,
				})
				gen = append(gen, res.Gen...)
				empty = append(empty, res.Empty...)
			}
			tc, ok := testcases[rel]
			if !ok {
				t.Fatalf("invalid test case: %s", rel)
			}
			if len(empty) > 0 {
				t.Errorf("got %d empty rules; want 0", len(empty))
			}
			f := rule.EmptyFile("test", "")
			for _, r := range gen {
				r.Insert(f)
			}
			for _, l := range langs {
				loads = append(loads, l.(language.ModuleAwareLanguage).ApparentLoads(func(string) string { return "" })...)
			}

			merger.FixLoads(f, loads)
			f.Sync()
			got := string(bzl.Format(f.File))
			wantPath := filepath.Join(dir, "BUILD.want")
			wantBytes, err := os.ReadFile(wantPath)
			if err != nil {
				t.Fatalf("error reading %s: %v", wantPath, err)
			}
			want := string(wantBytes)
			want = strings.ReplaceAll(want, "\r\n", "\n")
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("(-want, +got): %s", diff)
			}
			for kind, count := range tc.RulesCount {
				assert.Equal(t, len(filterRulesByKind(gen, kind)), count, tc.ErrorMessage)
			}
		})
	})
}

func filterRulesByKind(rules []*rule.Rule, name string) []*rule.Rule {
	var filteredRules []*rule.Rule
	for _, r := range rules {
		if r.Kind() == name {
			filteredRules = append(filteredRules, r)
		}
	}
	return filteredRules
}
