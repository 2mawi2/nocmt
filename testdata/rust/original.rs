#![allow(unused_variables)]
#![feature(async_closure)]

// First comment line
/* Header block */

use std::fmt; // inline module comment

// Empty comment lines
//
//
//
//

#[derive(Debug)]
struct Point {
    x: i32, /* inline */
    // y coordinate
    y: i32,
}

/* Nested-style
   /* inner */
   outer end */
fn compute() -> i32 {
    //goofy comment with code fn fake() {}
    42 /* trailing */
}

// Attribute that must stay
#[cfg(test)]
mod tests {
    // Test comment
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4); // trailing
    }
} 