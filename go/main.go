package main

import (
	"fmt"
	"os"
)

func main() {
	c := make([]byte, 1)
	for {
		os.Stdin.Read(c)
		if string(c) == "q" {
			break
		}
		fmt.Printf("Press your key is %v, '%s'\n", c, string(c))
	}
}
