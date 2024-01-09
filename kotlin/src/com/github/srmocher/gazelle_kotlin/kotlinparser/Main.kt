package com.github.srmocher.gazelle_kotlin.kotlinparser // ktlint-disable package-name

class Main {
    companion object {
        @JvmStatic fun main(args: Array<String>) {
            val port = System.getenv("KOTLIN_PARSER_PORT")?.toInt() ?: 50051
            val server = KotlinParserServer(port)
            server.start()
            server.blockUntilShutdown()
        }
    }
}
