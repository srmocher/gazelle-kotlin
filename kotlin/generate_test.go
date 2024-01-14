package kotlin

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
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

func testConfig(t *testing.T, repoRoot string) (*config.Config, language.Language, []config.Configurer) {
	cexts := []config.Configurer{
		&config.CommonConfigurer{},
		&walk.Configurer{},
		&resolve.Configurer{},
	}
	lang := NewLanguage()
	c := testtools.NewTestConfig(t, cexts, []language.Language{lang}, []string{
		"-build_file_name=BUILD.old",
		"-repo_root=" + repoRoot,
	})
	cexts = append(cexts, lang)
	return c, lang, cexts
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
	c, lang, cexts := testConfig(t, "testdata")
	content := []byte(`
# gazelle:follow **
# gazelle:kt_enabled true
# gazelle:kt_maven_repo maven_install.json maven
`)
	f, err := rule.LoadData(filepath.FromSlash("BUILD.config"), "config", content)
	if err != nil {
		t.Fatal(err)
	}
	for _, cext := range cexts {
		cext.Configure(c, "", f)
	}

	finishable := lang.(language.FinishableLanguage)
	defer finishable.DoneGeneratingRules()

	walk.Walk(c, []config.Configurer{lang}, []string{"testdata"}, walk.VisitAllUpdateSubdirsMode, func(dir, rel string, c *config.Config, update bool, oldFile *rule.File, subdirs, regularFiles, genFiles []string) {
		if !isTest(regularFiles) {
			return
		}
		t.Run(rel, func(t *testing.T) {
			res := lang.GenerateRules(language.GenerateArgs{
				Config:       c,
				Dir:          dir,
				Rel:          rel,
				File:         oldFile,
				Subdirs:      subdirs,
				RegularFiles: regularFiles,
				GenFiles:     genFiles,
			})
			if len(res.Empty) > 0 {
				t.Errorf("got %d empty rules; want 0", len(res.Empty))
			}
			f := rule.EmptyFile("test", "")
			for _, r := range res.Gen {
				r.Insert(f)
			}
			merger.FixLoads(f, lang.Loads())
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
			assert.NotEmptyf(t, res.Gen, "got %d generated rules; want > 0", len(res.Gen))

		})
	})
}
