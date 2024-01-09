package com.github.srmocher.gazelle_kotlin.kotlinparser.testdata

import kotlinx.coroutines.runBlocking
import java.time.LocalDate
// another simple example
class ExampleFile2 {
    companion object {
        val helloWorld = "hello world"
    }
    fun main(args: Array<String>) {
        runBlocking {
            val dt = LocalDate.now()

            println("The date is: $dt")
        }

        println("This will run last")
    }
}
