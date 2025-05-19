//go:build linux && !windows
// +build linux,!windows

/*
Package documentation

	spanning multiple lines
*/
package main

import (
	"fmt"
)

const Version = "v1.0.0"

/* Block comment before var */
var name = /* mid-token */ "Gopher"
var age = 10 /* trailing block */

func hello() {
	fmt.Println("Hello") /* inline block comment */
}

/*
 * Multi-line block
 * comment inside file
 */
func main() {
	hello()

	//go:generate echo "generate something"
	//go:noinline

	if true {
		fmt.Println("Conditional")
	} else /* else comment */ {
		fmt.Println("Else branch")
	}

}
