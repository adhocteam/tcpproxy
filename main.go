package main

import (
	"flag"
	"io"
	"log"
	"net"
)

func main() {
	var (
		listenAddr = flag.String("l", "", "local address to listen on")
		remoteAddr = flag.String("r", "", "remote address to dial")
	)

	flag.Parse()

	if *listenAddr == "" {
		log.Fatalf("must supply local address to listen on, -l option")
	}

	if *remoteAddr == "" {
		log.Fatalf("must supply remote address to dial, -r option")
	}

	ln, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalf("listening: %v", err)
	}

	proxy(ln, *remoteAddr)
}

func proxy(ln net.Listener, remoteAddr string) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		log.Printf("connected: %s", conn.RemoteAddr())

		go handle(conn, remoteAddr)
	}
}

func handle(conn net.Conn, remoteAddr string) {
	defer conn.Close()

	rconn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("dialing remote: %v", err)
		return
	}
	defer rconn.Close()

	copy(conn, rconn)
}

func copy(a, b io.ReadWriter) {
	done := make(chan struct{})
	go func() {
		io.Copy(a, b)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(b, a)
		done <- struct{}{}
	}()
	<-done
	<-done
}
