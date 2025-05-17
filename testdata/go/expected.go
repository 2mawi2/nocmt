//go:build linux && !windows
// +build linux,!windows

package main

import (
	"fmt" 
)

const Version = "v1.0.0" 

var (

	name =  "Gopher"

	age = 10 
)

func hello() { 

	fmt.Println("Hello")  
}

func main() {
	hello()

	//go:generate echo "generate something"

	//go:noinline

	if true { 
		fmt.Println("Conditional")
	} else  {
		fmt.Println("Else branch")
	}

}
