/*
 * @license
 * ExampleService.java
 * Copyright 2025 Example Corp
 */

/**
 * This is a simple service class for demonstration purposes.
 */
public class ExampleService {

    // Logger instance for this service
    private static final java.util.logging.Logger logger = java.util.logging.Logger.getLogger(ExampleService.class.getName());

    // Service name
    private final String serviceName;

    /**
     * Constructs a new ExampleService with the given name.
     * @param serviceName the name of the service
     */
    public ExampleService(String serviceName) {
        this.serviceName = serviceName;
        logger.info("ExampleService created with name: " + serviceName);
    }

    /**
     * Greets the user with a welcome message.
     * @param user the user's name
     * @return the greeting message
     */
    public String greetUser(String user) {
        // Compose the greeting message
        String message = "Hello, " + user + "! Welcome to " + serviceName + ".";
        logger.fine("Greeting generated: " + message); // Log at FINE level
        return message;
    }

    /**
     * Shuts down the service.
     */
    public void shutdown() {
        // Perform shutdown logic here
        logger.warning("Service " + serviceName + " is shutting down."); // Log at WARNING level
    }

    public static void main(String[] args) {
        // Create the service
        ExampleService service = new ExampleService("DemoService");

        // Greet a few users
        String[] users = {"Alice", "Bob", "Charlie"};
        for (String user : users) {
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