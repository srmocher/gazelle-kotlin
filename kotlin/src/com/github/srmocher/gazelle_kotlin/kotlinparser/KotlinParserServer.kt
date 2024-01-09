package com.github.srmocher.gazelle_kotlin.kotlinparser // ktlint-disable package-name

import arrow.core.Either
import io.grpc.Server
import io.grpc.ServerBuilder
import io.grpc.Status
import io.grpc.StatusException

data class KotlinParserError(val message: String)
class KotlinParserServer(private val port: Int) {
    val server: Server = ServerBuilder.forPort(port)
        .addService(KotlinParserService())
        .build()

    fun start() {
        server.start()
        println("Server started, listening on $port")
        Runtime.getRuntime().addShutdownHook(
            Thread {
                println("*** shutting down gRPC server since JVM is shutting down")
                this@KotlinParserServer.stop()
                println("*** server shut down")
            },
        )
    }

    private fun stop() {
        server.shutdown()
    }

    fun blockUntilShutdown() {
        server.awaitTermination()
    }

    internal class KotlinParserService : KotlinParserGrpcKt.KotlinParserCoroutineImplBase() {
        val parser = Parser()
        override suspend fun parseKotlinFiles(request: KotlinParserRequest): KotlinParserResponse {
            val result = parser.parseKtFiles(request.kotlinSourceFileList)
            when (result) {
                is Either.Left -> throw StatusException(Status.INTERNAL)
                is Either.Right -> {
                    val response = KotlinParserResponse.newBuilder()
                        .addAllSourceFileInfos(result.value)
                        .build()
                    return response
                }
            }
        }

        override suspend fun parseJavaFiles(request: JavaParserRequest): JavaParserResponse {
            return javaParserResponse {
                sourceFileInfo {
                    fileName = "dummy.java"
                    packageName = "package2"
                }
            }
        }
    }
}
