/*
 * @license
 * ExampleService.java
 * Copyright 2025 Example Corp
 */

package com.example.service;

import java.util.logging.Logger;
import java.util.List;
import java.util.ArrayList;

/**
 * This is a simple service class for demonstration purposes.
 * @author John Doe
 * @version 1.0
 */
@Component
@Service
public class ExampleService {

    private static final Logger logger = Logger.getLogger(ExampleService.class.getName());

    private final String serviceName;

    /**
     * Constructs a new ExampleService with the given name.
     * @param serviceName the name of the service
     */
    @Inject
    public ExampleService(String serviceName) {
        this.serviceName = serviceName;
        logger.info("ExampleService created with name: " + serviceName);
    }

    /**
     * Greets the user with a welcome message.
     * @param user the user's name
     * @return the greeting message
     */
    @Override
    @SuppressWarnings("unchecked")
    public String greetUser(String user) {
        String message = "Hello, " + user + "! Welcome to " + serviceName + ".";
        logger.fine("Greeting generated: " + message);
        return message;
    }

    /**
     * Processes a list of users.
     * @param users list of user names
     * @return processed results
     */
    @Deprecated("Use processUsersV2 instead")
    public List<String> processUsers(List<String> users) {
        List<String> results = new ArrayList<>();
        for (String user : users) {
            results.add(greetUser(user));
        }
        return results;
    }

    /**
     * Shuts down the service.
     */
    @PreDestroy
    public void shutdown() {
        logger.warning("Service " + serviceName + " is shutting down.");
    }

    public static void main(String[] args) {
        ExampleService service = new ExampleService("DemoService");

        String[] users = {"Alice", "Bob", "Charlie"};
        for (String user : users) {
            System.out.println(service.greetUser(user));
        }

        service.shutdown();
    }
}

/*
 * End of ExampleService.java
 */