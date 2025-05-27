<?php
#!/usr/bin/env php

# This is a shell-style comment
// This is a single-line comment
/* This is a multi-line comment */
/** This is a DocBlock comment */

class TestClass {
    // Constructor comment
    public function __construct() {
        # Another shell-style comment
        echo "Hello World"; // Inline comment
    }
    
    /* Multi-line comment
       spanning multiple lines */
    public function test() {
        /** 
         * DocBlock comment
         * with multiple lines
         */
        return true;
    }
}

// @license MIT
/* @preserve This should be preserved */
# @codingStandardsIgnoreStart
$variable = "value"; // End of line comment
# @codingStandardsIgnoreEnd 