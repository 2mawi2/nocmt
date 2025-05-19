/*
 * @license
 * ExampleService.java
 * Copyright 2025 Example Corp
 */

/**
 * This is a simple service class for demonstration purposes.
 */
public class ExampleService {

    private static final java.util.logging.Logger logger = java.util.logging.Logger.getLogger(ExampleService.class.getName());

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
        String message = "Hello, " + user + "! Welcome to " + serviceName + ".";
        logger.fine("Greeting generated: " + message);
        return message;
    }

    /**
     * Shuts down the service.
     */
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