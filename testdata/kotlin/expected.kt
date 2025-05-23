package com.example

import kotlin.collections.*

/**
 * Class documentation comment
 * @property name The name property
 */
@Entity
@Suppress("UNCHECKED_CAST")
class Person(
    private val name: String,
    val age: Int
) {
    private val id: Long = 0L
    
    /**
     * Function documentation 
     * @param greeting The greeting message
     * @return formatted greeting
     */
    @JvmStatic
    fun greet(greeting: String): String {
        return "$greeting, $name!"
    }
    
    /* 
     * Multi-line comment block
     * with multiple lines
     */
    private fun validate() {
        if (name.isEmpty()) {
            throw IllegalArgumentException("Name cannot be empty")
        }
    }
}

@Deprecated("Use factory method instead")
object PersonFactory {
    fun createPerson(name: String, age: Int): Person {
        /* Block comment before return */
        return Person(name, age)
    }
} 