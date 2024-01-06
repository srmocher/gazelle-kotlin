package kotlin

import (
    "github.com/bazelbuild/bazel-gazelle/rule"
)

var kotlinKinds = map[string]rule.KindInfo{
    "kt_jvm_library":  {
        MatchAny: true,
        NonEmptyAttrs: map[string]bool{
			"deps":  true,
			"srcs":  true,
		},
        MergeableAttrs: map[string]bool{
			"srcs": true,
            "exports": true,
            "javac_opts": true,
            "kotlinc_opts": true,
			"runtime_deps": true,
			"resource_strip_prefix": true,
		},
    }
}