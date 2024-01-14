# gazelle-kotlin
![CI](https://github.com/srmocher/gazelle-kotlin/actions/workflows/ci.yaml/badge.svg)
This is currently a early work-in-progress implementation of a [Gazelle](https://github.com/bazelbuild/bazel-gazelle) extension for Kotlin to generate BUILD targets. It uses rules from [rules_kotlin](https://github.com/bazelbuild/rules_kotlin) and [grpc-kotlin](https://github.com/grpc/grpc-kotlin) to generate targets within the source tree and relies on [rules_jvm_external](https://github.com/bazelbuild/rules_jvm_external) to retrieve information of external dependencies.

## Usage

This implementation is inspired by the Java extension in [rules_jvm](https://github.com/bazel-contrib/rules_jvm/tree/main/java/gazelle). However, it is currently implemented as a standalone extension and does not depend on that repository.

### WORKSPACE

TODO: once a release is published and `repositories.bzl` is setup.

### Bzlmod

TODO: Add bzlmod support

### Using the extension
To enable the new extension, add

```python
# gazelle:lang kotlin
# gazelle:kt_enabled true
# gazelle:kt_generate_proto true
# gazelle:kt_maven_repo maven_install.json maven
```

to your `BUILD.bazel` file. This will enable the extension for the current package and all subpackages. The `kt_enabled` flag is required to enable the extension. The `kt_generate_proto` flag is required to generate proto targets. The `kt_maven_repo` directive is required to extract external dependencies information and the value of it must be the relative path to the maven lock file followed by the maven Bazel repository name (typically `maven`).

Supported features:
- Generating `kt_jvm_library` targets for packages containing `.kt` files.
- Generating `kt_jvm_proto_library` targets for packages containing `.proto` files.
- Generating `kt_jvm_grpc_library` targets for packages containing `.proto` files with gRPC services.

Features to be added:
- Adding support for `kt_jvm_test`
- Adding support for `resources` and Java source files.
- 