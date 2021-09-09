package main

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func enableRawMode() {
	termios, err := unix.IoctlGetTermios(syscall.Stdin, unix.TIOCGETA)
	if err != nil {
		panic(err)
	}

	termios.Lflag &^= unix.ECHO

	//TODO TIOCSETA is maybe system dependent code? You can use TCSETA?
	if err := unix.IoctlSetTermios(syscall.Stdin, unix.TIOCSETA, termios); err != nil {
		panic(err)
	}
}

func main() {
	enableRawMode()

	c := make([]byte, 1)
	for {
		os.Stdin.Read(c)
		if string(c) == "q" {
			break
		}
		fmt.Printf("Press your key is %v, '%s'\n", c, string(c))
	}
}
