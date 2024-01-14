package com.foo

import java.time.LocalDateTime
import com.bar.*

fun main(args: Array<String>) {
    val current = LocalDateTime.now()
    var bar = Bar.newBuilder()
            .setMessage("hello world!")
            .build()

    println("Current Date and Time is: $current")
    println("bar says ${bar.message}")
}
