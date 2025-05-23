/*
 * Comprehensive Java test file
 * Tests various comment scenarios and Java features
 */

package com.example.test;

// Import statement comment
import java.util.*;
import java.util.stream.Collectors; // Stream import comment

/**
 * Comprehensive test class for Java comment removal
 * This javadoc should be preserved
 * @author Test Author
 * @version 1.0
 */
@SuppressWarnings("unchecked") // Suppress warnings annotation
@Deprecated("Use newer implementation") // Deprecation annotation
public class ComprehensiveTest<T extends Comparable<T>> { // Generic class comment

    // @formatter:off
    private static final String CONSTANT = "value"; // Formatter directive should be preserved
    // @formatter:on

    // Private field comment
    private final List<T> items;
    
    /* Multi-line comment
     * spanning multiple lines
     * should be preserved */
    private Map<String, Integer> cache = new HashMap<>(); // Inline comment to remove

    /**
     * Constructor documentation
     * @param items Initial items list
     */
    public ComprehensiveTest(List<T> items) {
        // Constructor body comment
        this.items = new ArrayList<>(items);
    }

    // CHECKSTYLE:OFF - disable checkstyle
    @Override // Override annotation comment
    public String toString() {
        // Method implementation comment
        return items.stream() // Stream comment
            .map(Object::toString) // Mapping comment
            .collect(Collectors.joining(", ")); // Collection comment
    }
    // CHECKSTYLE:ON - enable checkstyle

    /**
     * Generic method with lambda
     * @param <R> Return type
     * @param mapper Function to apply
     * @return Mapped results
     */
    public <R> List<R> map(java.util.function.Function<T, R> mapper) {
        // Lambda implementation comment
        return items.stream()
            .map(mapper) // Apply mapper function
            .collect(Collectors.toList()); // NOSONAR - collect to list
    }

    // @SuppressWarnings("rawtypes") - directive comment
    @SuppressWarnings("rawtypes")
    private void rawTypeMethod() {
        // Raw type usage comment
        List rawList = new ArrayList(); // Raw list comment
        rawList.add("item"); // Add item comment
    }

    /**
     * Method with try-catch block
     */
    public void methodWithException() {
        try {
            // Try block comment
            String result = performOperation(); // Operation call comment
            System.out.println(result); // Print result comment
        } catch (Exception e) {
            // Catch block comment
            e.printStackTrace(); // Print stack trace comment
        } finally {
            // Finally block comment
            cleanup(); // Cleanup call comment
        }
    }

    // Private helper method comment
    private String performOperation() throws Exception {
        // Method body comment
        if (items.isEmpty()) { // Check if empty comment
            throw new IllegalStateException("No items available"); // Exception comment
        }
        return "Operation completed"; // Return comment
    }

    // Another helper comment  
    private void cleanup() {
        // Cleanup implementation comment
        cache.clear(); // Clear cache comment
    }

    // Static nested class comment
    public static class NestedClass {
        // Nested class field comment
        private String value;

        // Nested class constructor comment
        public NestedClass(String value) {
            this.value = value; // Assignment comment
        }

        // NOCHECKSTYLE - ignore style rules
        public String getValue() { return value; } // Getter comment
    }

    // Anonymous class example comment
    private Runnable createRunnable() {
        // Return anonymous class comment
        return new Runnable() {
            @Override // Anonymous override comment
            public void run() {
                // Anonymous method comment
                System.out.println("Running..."); // Print statement comment
            }
        }; // End anonymous class comment
    }

    // Main method comment
    public static void main(String[] args) {
        // Main method body comment
        List<String> testItems = Arrays.asList("a", "b", "c"); // Test data comment
        
        // Create instance comment
        ComprehensiveTest<String> test = new ComprehensiveTest<>(testItems);
        
        // Method calls comment
        System.out.println(test.toString()); // Print test comment
        
        // Lambda usage comment
        List<String> mapped = test.map(s -> s.toUpperCase()); // Map to uppercase comment
        System.out.println(mapped); // Print mapped comment
    }
}

// End of file comment 