package com.github.srmocher.gazelle_kotlin.kotlinparser

import arrow.core.Either
import arrow.core.left
import arrow.core.right
import com.google.devtools.build.runfiles.Runfiles
import org.junit.Assert.*
import org.junit.Test
import java.io.File
import java.lang.Exception

class RunfilesHelper {
    val runfiles = Runfiles.preload()
    fun getRunfilePath(workspacePath: String): String {
        return runfiles.withSourceRepository("gazelle-kotlin")
            .rlocation(workspacePath)
    }
}
internal class ParserTest {
    @Test
    fun testParsingSingleFileWorks() {
        val parser = Parser()
        val runfilesHelper = RunfilesHelper()
        val ktFile = runfilesHelper.getRunfilePath("gazelle-kotlin/kotlin/test/com/github/srmocher/gazelle_kotlin/kotlinparser/testdata/ExampleFile.kt")
        assertTrue(File(ktFile).exists())
        val result = parser.parseKtFile(ktFile)
        when (result) {
            is Either.Right -> {
                assertEquals(result.value.importsList.count(), 1)
                assertEquals(result.value.importsList[0], "java.time.LocalDate")
                assertEquals(result.value.packageName, "com.github.srmocher.gazelle_kotlin.kotlinparser.testdata")
            }
            is Either.Left -> fail("Parsing should work but failed with error ${result.value.message}!")
        }
    }

    @Test
    fun testParsingSingleFileWithAnnotationsWorks() {
        val parser = Parser()
        val runfilesHelper = RunfilesHelper()
        val ktFile = runfilesHelper.getRunfilePath("gazelle-kotlin/kotlin/test/com/github/srmocher/gazelle_kotlin/kotlinparser/testdata/ExampleFile2.kt")
        assertTrue(File(ktFile).exists())
        val result = parser.parseKtFile(ktFile)
        when (result) {
            is Either.Right -> {
                assertEquals(result.value.importsList.count(), 2)
                assertEquals(result.value.importsList[0], "kotlinx.coroutines.runBlocking")
                assertEquals(result.value.importsList[1], "java.time.LocalDate")
                assertEquals(result.value.packageName, "com.github.srmocher.gazelle_kotlin.kotlinparser.testdata")
            }
            is Either.Left -> fail("Parsing should work but failed with error ${result.value.message}!")
        }
    }

    @Test
    fun testParseInvalidFileFails() {
        val parser = Parser()
        val runfilesHelper = RunfilesHelper()
        val ktFile = runfilesHelper.getRunfilePath("gazelle-kotlin/kotlin/test/com/github/srmocher/gazelle_kotlin/kotlinparser/testdata/InvalidKotlinFile.kt")
        assertTrue(File(ktFile).exists())
        val result = parser.parseKtFile(ktFile)
        assertTrue(result.isLeft())
        assertNotNull(result.leftOrNull())
    }
}
