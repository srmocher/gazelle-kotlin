package kotlin

import "github.com/bazelbuild/bazel-gazelle/language"

const kotlinName = "kotlin"

type kotlinLang struct {
	// kotlinPkgRels is a set of relative paths to directories containing buildable
	// Go code. If the value is false, it means the directory does not contain
	// buildable Go code, but it has a subdir which does.
	kotlinPkgRels map[string]bool
}

func (*kotlinLang) Name() string { return kotlinName }

func NewLanguage() language.Language {
	return &kotlinLang{kotlinPkgRels: make(map[string]bool)}
}