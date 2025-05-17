@available(iOS 13.0, *)
// MARK: - Example
/* File header */

import Foundation // inline

// Empty lines
//
//
//
//

func greet() {  // trailing
    // inside comment
    print("Hello")  /* inline */
}

/* Multiline string test */
let message = """
Line 1
// not a comment
/* not a block */
"""

// Compiler directives to keep
#if DEBUG
// debug comment
func debugPrint() {
    print("Debug")
}
#endif 