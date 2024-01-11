package com.github.srmocher.gazelle_kotlin.kotlinparser

import java.net.ServerSocket


 // ktlint-disable package-name

class Main {
    companion object {
        @JvmStatic fun main(args: Array<String>) {
            val s = ServerSocket(0)
            val server = KotlinParserServer(s.localPort)
            server.start()
            server.blockUntilShutdown()
        }
    }
}
