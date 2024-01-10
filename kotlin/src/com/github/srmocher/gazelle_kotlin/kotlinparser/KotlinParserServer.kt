package com.github.srmocher.gazelle_kotlin.kotlinparser // ktlint-disable package-name

import arrow.core.Either
import com.google.rpc.Status
import com.google.rpc.StatusProto
import io.grpc.StatusException
import io.grpc.netty.NettyServerBuilder
import io.grpc.protobuf.services.ProtoReflectionService
import java.net.InetSocketAddress


data class KotlinParserError(val message: String)
class KotlinParserServer(private val port: Int) {
    val server = NettyServerBuilder.forAddress(InetSocketAddress("[::1]", port))
        .addService(KotlinParserService())
            .addService(ProtoReflectionService.newInstance())
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
            println("Received request to parse ${request.kotlinSourceFileList.count()}")
            return when (val result = parser.parseKtFiles(request.kotlinSourceFileList)) {
                is Either.Left -> {
                    val status = Status.newBuilder()
                            .setCode(13)
                            .setMessage(result.value.message)
                            .build()
                    KotlinParserResponse.newBuilder()
                            .setError(status)
                            .build()
                }

                is Either.Right -> {
                    KotlinParserResponse.newBuilder()
                            .addAllSourceFileInfos(result.value)
                            .build()
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
