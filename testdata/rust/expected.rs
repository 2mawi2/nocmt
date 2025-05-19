#![allow(unused_variables)]
#![feature(async_closure)]

/* Header block */

use std::fmt; // inline module comment


#[derive(Debug)]
struct Point {
    x: i32, /* inline */
    y: i32,
}

/* Nested-style
   /* inner */
   outer end */
fn compute() -> i32 {
    42 /* trailing */
}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4); // trailing
    }
} 