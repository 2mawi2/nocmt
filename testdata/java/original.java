/*
 * @license
 * ExampleService.java
 * Copyright 2025 Example Corp
 */

// Package declaration
package com.example.service;

// Import statements
import java.util.logging.Logger;
import java.util.List;
import java.util.ArrayList;

/**
 * This is a simple service class for demonstration purposes.
 * @author John Doe
 * @version 1.0
 */
@Component
@Service // Spring annotation
public class ExampleService {

    // Logger instance for this service
    private static final Logger logger = Logger.getLogger(ExampleService.class.getName());

    // Service name field
    private final String serviceName;

    /**
     * Constructs a new ExampleService with the given name.
     * @param serviceName the name of the service
     */
    @Inject // Dependency injection annotation
    public ExampleService(String serviceName) {
        this.serviceName = serviceName; // Assign service name
        logger.info("ExampleService created with name: " + serviceName);
    }

    /**
     * Greets the user with a welcome message.
     * @param user the user's name
     * @return the greeting message
     */
    @Override // Override annotation
    @SuppressWarnings("unchecked") // Suppress warnings
    public String greetUser(String user) {
        // Compose the greeting message
        String message = "Hello, " + user + "! Welcome to " + serviceName + ".";
        logger.fine("Greeting generated: " + message); // Log at FINE level
        return message; // Return the greeting
    }

    /**
     * Processes a list of users.
     * @param users list of user names
     * @return processed results
     */
    @Deprecated("Use processUsersV2 instead") // Deprecation annotation
    public List<String> processUsers(List<String> users) {
        // Initialize result list
        List<String> results = new ArrayList<>();
        for (String user : users) { // Iterate through users
            results.add(greetUser(user)); // Add greeting to results
        }
        return results; // Return processed results
    }

    /**
     * Shuts down the service.
     */
    @PreDestroy // Lifecycle annotation
    public void shutdown() {
        // Perform shutdown logic here
        logger.warning("Service " + serviceName + " is shutting down."); // Log at WARNING level
    }

    // Main method for testing
    public static void main(String[] args) {
        // Create the service
        ExampleService service = new ExampleService("DemoService");

        // Greet a few users
        String[] users = {"Alice", "Bob", "Charlie"};
        for (String user : users) { // Loop through users
            System.out.println(service.greetUser(user)); // Print greeting
        }

        // Simulate some operation
        // TODO: Add real business logic here

        // Shutdown the service
        service.shutdown();
    }
}

/*
 * End of ExampleService.java
 */