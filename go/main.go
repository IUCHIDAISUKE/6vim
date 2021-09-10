package main

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unicode"

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
	termios.Lflag &^= (unix.ECHO | unix.ICANON | unix.ISIG)

	//TODO TIOCSETA is maybe system dependent code? You can use TCSETA?
	if err := unix.IoctlSetTermios(syscall.Stdin, unix.TIOCSETA, termios); err != nil {
		panic(err)
	}
}

func main() {
	enableRawMode()
	defer disableRawMode()

	c := make([]byte, 4)
	for {
		n, err := os.Stdin.Read(c)
		if err != nil && err != io.EOF {
			panic(err)
		}
		r := rune(c[0])
		if unicode.IsControl(rune(c[0])) {
			fmt.Printf("%d\n", r)
		} else {
			fmt.Printf("%d ('%[1]c')\n", r)
		}
		if string(c[:n]) == "q" {
			break
		}
		c = []byte{0, 0, 0, 0}
	}
}
