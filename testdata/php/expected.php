<?php
#!/usr/bin/env php

/* This is a multi-line comment */
/** This is a DocBlock comment */

class TestClass {
    public function __construct() {
        echo "Hello World";
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
$variable = "value";
# @codingStandardsIgnoreEnd 