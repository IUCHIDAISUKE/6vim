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
	termios.Iflag &^= (unix.BRKINT | unix.ICRNL | unix.INPCK | unix.ISTRIP | unix.IXON)
	termios.Oflag &^= (unix.OPOST)
	termios.Cflag |= (unix.CS8)
	termios.Lflag &^= (unix.ECHO | unix.ICANON | unix.ISIG | unix.IEXTEN)
	termios.Cc[unix.VMIN] = 0
	termios.Cc[unix.VTIME] = 1

	//TODO TIOCSETA is maybe system dependent code? You can use TCSETA?
	if err := unix.IoctlSetTermios(syscall.Stdin, unix.TIOCSETA, termios); err != nil {
		panic(err)
	}
}

func main() {
	enableRawMode()
	defer disableRawMode()

	for {
		//TODO too many allocate memory?
		c := make([]byte, 4)

		n, err := os.Stdin.Read(c)
		if err != nil && err != io.EOF {
			panic(err)
		}
		r := rune(c[0])
		if unicode.IsControl(rune(c[0])) {
			fmt.Printf("%d\r\n", r)
		} else {
			fmt.Printf("%d ('%[1]c')\r\n", r)
		}
		if string(c[:n]) == "q" {
			break
		}
	}
}
