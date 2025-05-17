#![allow(unused_variables)]
#![feature(async_closure)]

use std::fmt; 

#[derive(Debug)]
struct Point {
    x: i32, 

    y: i32,
}

fn compute() -> i32 {

    42 
}

#[cfg(test)]
mod tests {

    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4); 
    }
} 
