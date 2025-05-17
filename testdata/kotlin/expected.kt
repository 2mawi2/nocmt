// @file:JvmName("SampleFile")
package example

import kotlin.time.*

class Greeter {

    val greeting = "Hi"

    // @Suppress("RedundantSuspendModifier")
    suspend fun greet() {
        println(greeting)
    }
}

// @OptIn(ExperimentalTime::class)
fun timing(block: () -> Unit) {
    measureTime {  block() }
} 