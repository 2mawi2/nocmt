// Copyright 2025 Example
// License: MIT

//go:build linux && !windows
// +build linux,!windows

/*
Package documentation

	spanning multiple lines
*/
package main

import (
	"fmt" // inline import comment
)

// Single-line comment before constant
const Version = "v1.0.0" // trailing comment after constant

/* Block comment before var */
var name = /* mid-token */ "Gopher"
var age = 10 /* trailing block */

// Empty comment lines
//
//
//
//

func hello() { // trailing comment after func sig
	// Regular line comment
	fmt.Println("Hello") /* inline block comment */ // twin trailing
}

/*
 * Multi-line block
 * comment inside file
 */
func main() {
	hello()

	//go:generate echo "generate something"
	// Regular comment between directives
	//go:noinline

	if true { // conditional comment
		fmt.Println("Conditional")
	} else /* else comment */ {
		fmt.Println("Else branch")
	}

	// URL-like comment
	// https://golang.org/pkg/fmt
}
