package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

func execute(conn net.Conn) {
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	go func() { _, _ = io.Copy(os.Stdin, conn) }()
	_, _ = io.Copy(conn, os.Stdout)

}

func main() {
	port := os.Args[1]
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	l, err := tls.Listen("tcp4", ":"+port, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	if err != nil {
		println(err.Error())
	}
	println("waiting")
	conn, _ := l.Accept()
	println("recvd")

	execute(conn)

	// do something to determine when to stop
}
