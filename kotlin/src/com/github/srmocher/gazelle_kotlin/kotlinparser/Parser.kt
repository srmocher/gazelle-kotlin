package com.github.srmocher.gazelle_kotlin.kotlinparser // ktlint-disable package-name

import arrow.core.*
import arrow.core.raise.catch
import arrow.core.raise.either
import arrow.fx.coroutines.parMap
import kotlinx.ast.common.AstSource
import kotlinx.ast.common.ast.Ast
import kotlinx.ast.common.ast.DefaultAstNode
import kotlinx.ast.common.klass.identifierName
import kotlinx.ast.grammar.kotlin.common.summary
import kotlinx.ast.grammar.kotlin.common.summary.PackageHeader
import kotlinx.ast.grammar.kotlin.target.antlr.kotlin.KotlinGrammarAntlrKotlinParser
import java.io.File
import java.lang.Exception

class Parser {
    private fun astsToSourceInfo(file: String, asts: List<Ast>): SourceFileInfo {
        var pkg: String = ""
        val importsList = mutableListOf<String>()
        asts.forEach { ast: Ast ->
            if (ast is PackageHeader) {
                pkg = ast.identifier.identifierName()
            } else if ((ast is DefaultAstNode) && (ast.description.equals("importList"))) {
                val impList = ast.children.map { impNode -> impNode.description.replace("Import(", "").replace(")", "") }
                importsList.addAll(impList)
            }
        }
        val si = SourceFileInfo.newBuilder()
            .addAllImports(importsList)
            .setFileName(file)
            .setPackageName(pkg)
            .build()
        return si
    }

    fun parseKtFile(file: String): Either<KotlinParserError, SourceFileInfo> {
        val f = File(file)
        if (!f.exists()) {
            return KotlinParserError("File $file does not exist").left()
        }
        val source = AstSource.File(file)
        catch({
            val parseResult = KotlinGrammarAntlrKotlinParser.parseKotlinFile(source)
            val summary = parseResult.summary(attachRawAst = true)
            return when (summary.errorList().count()) {
                0 -> astsToSourceInfo(file, summary.get()).right()
                else -> KotlinParserError("Error in parsing: ${summary.errorList().joinToString("\n")}").left()
            }
        }) { e: Exception ->
            KotlinParserError("Error parsing kotlin file $file: ${e.message}")
        }

        return KotlinParserError("Failed to parse for unknown reasons").left()
    }

    // Parses multiple Kotlin source files in parallel and returns the result
    // which is either an error in parsing any one of the files or list of SourceFileInfo
    // objects which indicate the metadata of each of the source files that were parsed
    suspend fun parseKtFiles(files: List<String>): Either<KotlinParserError, List<SourceFileInfo>> {
        return files.parMap { file -> parseKtFile(file) }.let { l -> either { l.bindAll() } }
    }
}
