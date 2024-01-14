package com.github.srmocher.gazelle_kotlin.kotlinparser // ktlint-disable package-name

import arrow.core.Either
import com.google.rpc.Status
import com.google.rpc.StatusProto
import io.grpc.Server
import io.grpc.StatusException
import io.grpc.netty.NettyServerBuilder
import io.grpc.protobuf.services.ProtoReflectionService
import java.net.InetSocketAddress
import java.nio.charset.StandardCharsets
import java.nio.file.Files
import java.nio.file.Paths


data class KotlinParserError(val message: String)
class KotlinParserServer(private val port: Int) {
    private val server: Server = NettyServerBuilder.forAddress(InetSocketAddress("[::1]", port))
        .addService(KotlinParserService())
            .addService(ProtoReflectionService.newInstance())
        .build()

    private val pid = ProcessHandle.current().pid()
    private val portFile = Paths.get("/tmp", "gazelle-kotlin", "kotlinparser.${pid}.port")

    private fun storePortFile() {
        Files.createDirectories(Paths.get("/tmp", "gazelle-kotlin"))
        Files.write(portFile, listOf(port.toString()), StandardCharsets.UTF_8)
    }

    fun start() {
        server.start()
        println("Server started, listening on $port")
        Runtime.getRuntime().addShutdownHook(
            Thread {
                println("*** shutting down gRPC server since JVM is shutting down")
                Files.deleteIfExists(portFile)
                this@KotlinParserServer.stop()
                println("*** server shut down")
            },
        )
        storePortFile()
    }

    private fun stop() {
        server.shutdown()
    }

    fun blockUntilShutdown() {
        server.awaitTermination()
    }

    internal class KotlinParserService : KotlinParserGrpcKt.KotlinParserCoroutineImplBase() {
        private val parser = Parser()
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
