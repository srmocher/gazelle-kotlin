package kotlin

import (
	"fmt"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

var kotlinKinds = map[string]rule.KindInfo{
	"kt_jvm_library": {
		MatchAny: true,
		NonEmptyAttrs: map[string]bool{
			"deps": true,
			"srcs": true,
		},
		MergeableAttrs: map[string]bool{
			"srcs":                  true,
			"exports":               true,
			"javac_opts":            true,
			"kotlinc_opts":          true,
			"runtime_deps":          true,
			"resource_strip_prefix": true,
		},
	},
}

var kotlinLoads = []rule.LoadInfo{
	{
		Name: "@io_bazel_rules_kotlin//kotlin:jvm.bzl",
		Symbols: []string{
			"kt_jvm_library",
			"kt_jvm_test",
			"kt_jvm_binary",
		},
	},
	{
		Name: "@com_github_grpc_grpc_kotlin//:kt_jvm_grpc.bzl",
		Symbols: []string{
			"kt_jvm_grpc_library",
			"kt_jvm_proto_library",
		},
	},
}

func (*kotlinLang) Loads() []rule.LoadInfo {
	return kotlinLoads
}

func (*kotlinLang) Kinds() map[string]rule.KindInfo {
	return kotlinKinds
}

func (*kotlinLang) ApparentLoads(moduleToApparentName func(string) string) []rule.LoadInfo {
	return apparentLoads(moduleToApparentName)
}

func apparentLoads(moduleToApparentName func(string) string) []rule.LoadInfo {
	rulesGo := moduleToApparentName("rules_kotlin")
	if rulesGo == "" {
		rulesGo = "io_bazel_rules_kotlin"
	}

	return []rule.LoadInfo{
		{
			Name: fmt.Sprintf("@%s//kotlin:jvm.bzl", rulesGo),
			Symbols: []string{
				"kt_jvm_library",
				"kt_jvm_test",
				"kt_jvm_binary",
			},
		}, {
			Name: fmt.Sprintf("@com_github_grpc_grpc_kotlin//:kt_jvm_grpc.bzl"),
			Symbols: []string{
				"kt_jvm_grpc_library",
				"kt_jvm_proto_library",
			},
		},
	}
}
