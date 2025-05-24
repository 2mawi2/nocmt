// Single line comment at the top
import Foundation

/// Documentation comment for the class
/// This is a multi-line documentation comment
class SwiftExample {
    // Property comment
    var name: String = "example" // Inline comment
    
    /* Multi-line comment
       spanning multiple lines
       with more details */
    private var count: Int = 0
    
    /**
     * JavaDoc style documentation comment
     * @param name The name parameter
     * @returns Void
     */
    func setName(_ name: String) {
        // Simple line comment
        self.name = name // Another inline comment
    }
    
    // TODO: Implement this method
    func increment() {
        /* Block comment */ count += 1
    }
    
    /// Single line doc comment
    func getValue() -> Int {
        return count // Return comment
    }
}

// Final comment at the end 