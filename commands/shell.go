package commands

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
)

type termCloser struct {
	once   sync.Once
	f      *os.File
	closed bool
}

func (c *termCloser) Close() (err error) {
	c.once.Do(func() {
		err = c.f.Close()
		c.closed = true
	})
	return
}

func execShell(command string, c *termCloser) {
	cmd := exec.Command("/bin/bash", "-c", command)
	if c == nil {
		c = new(termCloser)
	}

	ptmx, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	c.f = ptmx

	// Make sure to close the pty at the end.
	defer func() { _ = c.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	_, _ = io.Copy(os.Stdout, ptmx)
}

func execAndReturn(command string) string {
	cmd := exec.Command("/bin/bash", "-c", command)
	f, err := pty.Start(cmd)
	if err != nil {
		panic(err)
	}
	recOutput, _ := os.Create("/tmp/output")
	buf := new(strings.Builder)
	//if _, err := io.Copy(buf, f); err != nil {
	//	panic(err)
	//}
	io.Copy(buf, io.TeeReader(f, recOutput))
	return buf.String()
}
