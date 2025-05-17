// @file:JvmName("SampleFile")
package example

/* Header block comment */

import kotlin.time.* // inline

// 多行空注释
//
//
//
//

/**
 * Class doc
 */
class Greeter { // trailing
    // Property comment
    val greeting = "Hi" /* inline */

    // @Suppress("RedundantSuspendModifier")
    suspend fun greet() { // trailing
        println(greeting) // inline trailing
    }
}

// @OptIn(ExperimentalTime::class)
fun timing(block: () -> Unit) {
    measureTime { /* comment */ block() }
} 