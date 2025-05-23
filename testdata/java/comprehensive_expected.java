/*
 * Comprehensive Java test file
 * Tests various comment scenarios and Java features
 */

package com.example.test;

import java.util.*;
import java.util.stream.Collectors;

/**
 * Comprehensive test class for Java comment removal
 * This javadoc should be preserved
 * @author Test Author
 * @version 1.0
 */
@SuppressWarnings("unchecked")
@Deprecated("Use newer implementation")
public class ComprehensiveTest<T extends Comparable<T>> {

    // @formatter:off
    private static final String CONSTANT = "value";
    // @formatter:on

    private final List<T> items;
    
    /* Multi-line comment
     * spanning multiple lines
     * should be preserved */
    private Map<String, Integer> cache = new HashMap<>();

    /**
     * Constructor documentation
     * @param items Initial items list
     */
    public ComprehensiveTest(List<T> items) {
        this.items = new ArrayList<>(items);
    }

    // CHECKSTYLE:OFF - disable checkstyle
    @Override
    public String toString() {
        return items.stream()
            .map(Object::toString)
            .collect(Collectors.joining(", "));
    }
    // CHECKSTYLE:ON - enable checkstyle

    /**
     * Generic method with lambda
     * @param <R> Return type
     * @param mapper Function to apply
     * @return Mapped results
     */
    public <R> List<R> map(java.util.function.Function<T, R> mapper) {
        return items.stream()
            .map(mapper)
            .collect(Collectors.toList()); // NOSONAR - collect to list
    }

    // @SuppressWarnings("rawtypes") - directive comment
    @SuppressWarnings("rawtypes")
    private void rawTypeMethod() {
        List rawList = new ArrayList();
        rawList.add("item");
    }

    /**
     * Method with try-catch block
     */
    public void methodWithException() {
        try {
            String result = performOperation();
            System.out.println(result);
        } catch (Exception e) {
            e.printStackTrace();
        } finally {
            cleanup();
        }
    }

    private String performOperation() throws Exception {
        if (items.isEmpty()) {
            throw new IllegalStateException("No items available");
        }
        return "Operation completed";
    }

    private void cleanup() {
        cache.clear();
    }

    public static class NestedClass {
        private String value;

        public NestedClass(String value) {
            this.value = value;
        }

        // NOCHECKSTYLE - ignore style rules
        public String getValue() { return value; }
    }

    private Runnable createRunnable() {
        return new Runnable() {
            @Override
            public void run() {
                System.out.println("Running...");
            }
        };
    }

    public static void main(String[] args) {
        List<String> testItems = Arrays.asList("a", "b", "c");
        
        ComprehensiveTest<String> test = new ComprehensiveTest<>(testItems);
        
        System.out.println(test.toString());
        
        List<String> mapped = test.map(s -> s.toUpperCase());
        System.out.println(mapped);
    }
} 