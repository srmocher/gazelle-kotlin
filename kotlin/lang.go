package kotlin

import (
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/srmocher/gazelle-kotlin/kotlin/internal"
)

const kotlinName = "kotlin"

type kotlinLang struct {
	// kotlinPkgRels is a set of relative paths to directories containing buildable
	// Kotlin code. If the value is false, it means the directory does not contain
	// buildable Kotlin code, but it has a subdir which does.
	kotlinPkgRels map[string]bool

	// The parser handle used to get metadata about sources
	kotlinParser *internal.KotlinParser
}

func (*kotlinLang) Name() string { return kotlinName }

func NewLanguage() language.Language {
	return &kotlinLang{kotlinPkgRels: make(map[string]bool)}
}
