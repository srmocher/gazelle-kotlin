module github.com/srmocher/gazelle-kotlin

go 1.21.5

require (
	github.com/bazelbuild/bazel-gazelle v0.35.0
	github.com/bazelbuild/buildtools v0.0.0-20231115204819-d4c9dccdfbb1
	github.com/bazelbuild/rules_go v0.45.0
	github.com/google/go-cmp v0.6.0
	github.com/stretchr/testify v1.8.4
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231120223509-83a465c0220f
	google.golang.org/grpc v1.59.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/bmatcuk/doublestar/v4 v4.6.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools/go/vcs v0.1.0-deprecated // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/srmocher/gazelle-kotlin => ./
