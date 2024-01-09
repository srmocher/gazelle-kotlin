package com.github.srmocher.gazelle_kotlin.kotlinparser.testdata

import java.time.LocalDate
class ExampleFile {
    fun main(args: Array<String>) {
        val dt = LocalDate.now()

        println("The date is: $dt")
    }
}
