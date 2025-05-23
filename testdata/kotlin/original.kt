// File header comment
package com.example

// Import comments
import kotlin.collections.*

/**
 * Class documentation comment
 * @property name The name property
 */
@Entity
@Suppress("UNCHECKED_CAST") // Suppress warning annotation
class Person(
    // Constructor parameter comment
    private val name: String,
    val age: Int // Age property
) {
    // Property comment
    private val id: Long = 0L
    
    /**
     * Function documentation 
     * @param greeting The greeting message
     * @return formatted greeting
     */
    @JvmStatic // Platform annotation
    fun greet(greeting: String): String {
        // Single line comment
        return "$greeting, $name!" // Inline comment
    }
    
    /* 
     * Multi-line comment block
     * with multiple lines
     */
    private fun validate() {
        // Validation logic
        if (name.isEmpty()) { // Check name
            throw IllegalArgumentException("Name cannot be empty") // Error message
        }
    }
}

// Companion object comment
@Deprecated("Use factory method instead") // Deprecation annotation
object PersonFactory {
    // Factory method comment
    fun createPerson(name: String, age: Int): Person {
        /* Block comment before return */
        return Person(name, age)
    }
} 