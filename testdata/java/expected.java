package com.example.demo;

import java.util.ArrayList;
import java.util.List; 


public class Main {

    
    private static final String VERSION = "1.0.0"; 
    
    
    private List<String> items; 
    
    
    public Main() {
        
        this.items = new ArrayList<>();
    }
    
    // @formatter:off
    public void addItem(String item) {
        
        if (item != null && !item.isEmpty()) {
            this.items.add(item); 
        }
    }
    // @formatter:on
    
    
    public int processItems() {
        
        return this.items.size();
    }
    
    
    public static void main(String[] args) {
        
        Main app = new Main();
        
        // @SuppressWarnings("unused")
        int count = 0; 
        
        
        app.addItem("Item 1");
        app.addItem("Item 2"); 
        app.addItem("Item 3");
        
        
        System.out.println("Version: " + VERSION);
        
        
        count = app.processItems();
        
        
        System.out.println("Processed " + count + " items");
    }
} 