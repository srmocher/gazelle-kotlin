workspace(name = "gazelle-kotlin")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    integrity = "sha256-fHbWI2so/2laoozzX5XeMXqUcv0fsUrHl8m/aE8Js3w=",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.44.2/rules_go-v0.44.2.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.44.2/rules_go-v0.44.2.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    integrity = "sha256-MpOL2hbmcABjA1R5Bj2dJMYO2o15/Uc5Vj9Q0zHLMgk=",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.35.0/bazel-gazelle-v0.35.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.35.0/bazel-gazelle-v0.35.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

############################################################
# Define your own dependencies here using go_repository.
# Else, dependencies declared by rules_go/gazelle will be used.
# The first declaration of an external repository "wins".
############################################################

load("//:go_deps.bzl", "go_dependencies")

# gazelle:repository_macro go_deps.bzl%go_dependencies
go_dependencies()

go_rules_dependencies()

go_register_toolchains(version = "1.21.5")

gazelle_dependencies()

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

RULES_JVM_EXTERNAL_TAG = "4.5"

RULES_JVM_EXTERNAL_SHA = "b17d7388feb9bfa7f2fa09031b32707df529f26c91ab9e5d909eb1676badd9a6"

http_archive(
    name = "rules_jvm_external",
    sha256 = RULES_JVM_EXTERNAL_SHA,
    strip_prefix = "rules_jvm_external-%s" % RULES_JVM_EXTERNAL_TAG,
    url = "https://github.com/bazelbuild/rules_jvm_external/archive/%s.zip" % RULES_JVM_EXTERNAL_TAG,
)

load("@rules_jvm_external//:repositories.bzl", "rules_jvm_external_deps")

rules_jvm_external_deps()

load("@rules_jvm_external//:setup.bzl", "rules_jvm_external_setup")

rules_jvm_external_setup()

load("@rules_jvm_external//:defs.bzl", "maven_install")

GRPC_JAVA_TAG = "1.60.1"

GRPC_KOTLIN_TAG = "1.4.1"

http_archive(
    name = "io_grpc_grpc_java",
    sha256 = "fedf186a6a66f7aff0d9be8d194adcb02426130b9942de4ea314400b6b44d528",
    strip_prefix = "grpc-java-{version}".format(version = GRPC_JAVA_TAG),
    url = "https://github.com/grpc/grpc-java/archive/refs/tags/v{version}.tar.gz".format(version = GRPC_JAVA_TAG),
)

http_archive(
    name = "com_github_grpc_grpc_kotlin",
    sha256 = "ebe0497ce15bc299501edaa86a0cb9586cc3ac70bf8464a3c1c3df6615260223",
    strip_prefix = "grpc-kotlin-{version}".format(version = GRPC_KOTLIN_TAG),
    url = "https://github.com/grpc/grpc-kotlin/archive/refs/tags/v{version}.tar.gz".format(version = GRPC_KOTLIN_TAG),
)

load("@io_grpc_grpc_java//:repositories.bzl", "IO_GRPC_GRPC_JAVA_ARTIFACTS", "IO_GRPC_GRPC_JAVA_OVERRIDE_TARGETS", "grpc_java_repositories")

grpc_java_repositories()

load("@com_github_grpc_grpc_kotlin//:repositories.bzl", "IO_GRPC_GRPC_KOTLIN_ARTIFACTS", "IO_GRPC_GRPC_KOTLIN_OVERRIDE_TARGETS", "grpc_kt_repositories")

grpc_kt_repositories()

load("@com_google_protobuf//:protobuf_deps.bzl", "PROTOBUF_MAVEN_ARTIFACTS")
load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

maven_install(
    artifacts = [
        "io.grpc:grpc-api:1.40.0",
        "io.grpc:grpc-core:1.40.0",
        "io.grpc:grpc-netty:1.40.0",
        "io.grpc:grpc-services:1.40.0",
        "io.grpc:grpc-stub:1.40.0",
        "org.slf4j:slf4j-simple:1.7.32",
        "com.google.code.findbugs:jsr305:3.0.2",
        "com.google.protobuf:protobuf-kotlin:3.24.0",

        # Library to parse Kotlin sources
        "com.github.kotlinx.ast:grammar-kotlin-parser-antlr-kotlin-jvm:v0.1.0",
        "com.github.kotlinx.ast:grammar-kotlin-parser-common-jvm:v0.1.0",
        "com.github.kotlinx.ast:common-jvm:v0.1.0",
        "com.github.kotlinx.ast:parser-antlr-kotlin-jvm:v0.1.0",
        "io.arrow-kt:arrow-core-jvm:1.2.0",
        "io.arrow-kt:arrow-fx-coroutines-jvm:1.2.0",
        "org.jetbrains.kotlinx:kotlinx-coroutines-core-jvm:1.7.3",
        "org.jetbrains.kotlin:kotlin-test-junit5:1.8.20",
        "junit:junit:4.13.2",
    ] + IO_GRPC_GRPC_JAVA_ARTIFACTS + PROTOBUF_MAVEN_ARTIFACTS + IO_GRPC_GRPC_KOTLIN_ARTIFACTS,
    generate_compat_repositories = True,
    maven_install_json = "//:maven_install.json",
    override_targets = IO_GRPC_GRPC_JAVA_OVERRIDE_TARGETS,
    repositories = [
        "https://maven.google.com",
        "https://repo1.maven.org/maven2",
        "https://jitpack.io",
    ],
)

load("@maven//:defs.bzl", "pinned_maven_install")

pinned_maven_install()

load("@maven//:compat.bzl", "compat_repositories")

compat_repositories()

load("@io_bazel_rules_kotlin//kotlin:repositories.bzl", "kotlin_repositories")

kotlin_repositories()  # if you want the default. Otherwise see custom kotlinc distribution below

register_toolchains("//:kotlin_toolchain")  # to use the default toolchain, otherwise see toolchains below
