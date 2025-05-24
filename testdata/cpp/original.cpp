#include <iostream>
#include <vector>
#include <memory>

// TODO: Implement better error handling
// FIXME: Memory leak in constructor
// NOTE: This class is thread-safe

/* 
   This is a multi-line comment
   that should be preserved
*/

/**
 * Documentation comment for the class
 * @param value Initial value
 */
class Example {
private:
    int value; // Private member variable
    
public:
    // Constructor with initialization
    Example(int val) : value(val) {
        // Initialize the value
        std::cout << "Creating Example with value: " << val << std::endl;
    }
    
    // Getter method
    int getValue() const {
        return value; // Return the stored value
    }
    
    // Setter method  
    void setValue(int newVal) {
        value = newVal; // Update the value
    }
    
    /* Block comment in method */
    void processData() {
        // Local variable declaration
        std::vector<int> data = {1, 2, 3, 4, 5};
        
        // Process each element
        for (const auto& item : data) {
            // Print each item
            std::cout << item << " ";
        }
        std::cout << std::endl; // New line
    }
};

// Main function
int main() {
    // Create an instance
    Example ex(42);
    
    // Call methods
    ex.processData(); // Process the data
    
    // Get the value
    int val = ex.getValue(); // Store the value
    std::cout << "Value: " << val << std::endl;
    
    return 0; // Success
}

// End of file comment 