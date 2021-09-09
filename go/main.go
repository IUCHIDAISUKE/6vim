package main

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

var orig_termios *unix.Termios

func disableRawMode() {
	if err := unix.IoctlSetTermios(syscall.Stdin, unix.TIOCSETA, orig_termios); err != nil {
		panic(err)
	}
}

func enableRawMode() {
	var err error
	orig_termios, err = unix.IoctlGetTermios(syscall.Stdin, unix.TIOCGETA)
	if err != nil {
		panic(err)
	}

	termios := orig_termios
	termios.Lflag &^= unix.ECHO

	//TODO TIOCSETA is maybe system dependent code? You can use TCSETA?
	if err := unix.IoctlSetTermios(syscall.Stdin, unix.TIOCSETA, termios); err != nil {
		panic(err)
	}
}

func main() {
	enableRawMode()
	defer disableRawMode()

	c := make([]byte, 1)
	for {
		os.Stdin.Read(c)
		if string(c) == "q" {
			break
		}
		fmt.Printf("Press your key is %v, '%s'\n", c, string(c))
	}
}
