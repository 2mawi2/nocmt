import Foundation

/// Documentation comment for the class
/// This is a multi-line documentation comment
class SwiftExample {
    var name: String = "example"
    
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
        self.name = name
    }
    
    // TODO: Implement this method
    func increment() {
        /* Block comment */ count += 1
    }
    
    /// Single line doc comment
    func getValue() -> Int {
        return count
    }
} 