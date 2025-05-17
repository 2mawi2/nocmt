// Copyright 2025 Example Corp
// Licensed under MIT

package com.example.demo;

import java.util.ArrayList;
import java.util.List; // Import for lists

/**
 * Main application class
 * This is a sample program with comments
 *
 * @author Example Author
 */
public class Main {

    // Constants
    private static final String VERSION = "1.0.0"; /* Version number */
    
    // Class variables
    private List<String> items; // Items list
    
    /**
     * Constructor with documentation
     */
    public Main() {
        // Initialize items
        this.items = new ArrayList<>();
    }
    
    // @formatter:off
    public void addItem(String item) {
        // Check if item is valid
        if (item != null && !item.isEmpty()) {
            this.items.add(item); // Add to list
        }
    }
    // @formatter:on
    
    /*
     * Process all items in the list
     * and return the count
     */
    public int processItems() {
        // Return count
        return this.items.size();
    }
    
    /**
     * Main method
     * @param args Command line arguments
     */
    public static void main(String[] args) {
        // Create instance
        Main app = new Main();
        
        // @SuppressWarnings("unused")
        int count = 0; // Counter
        
        /* Add sample items */
        app.addItem("Item 1");
        app.addItem("Item 2"); // Second item
        app.addItem("Item 3");
        
        // Print version
        System.out.println("Version: " + VERSION);
        
        // Process items
        count = app.processItems();
        
        // Print result
        System.out.println("Processed " + count + " items");
    }
} 