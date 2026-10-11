package main

import (
	"flag"
	"io"
	"log"
	"net"
	"runtime"
	"time"
)

func main() {
	go func() {
		for {
			log.Printf("goroutines: %d", runtime.NumGoroutine())
			time.Sleep(time.Second)
		}
	}()

	// Flags listen and target
	listen := flag.String("listen", ":9000", "type a port for listening")
	target := flag.String("target", "localhost:6379", "target for conncetion")
	flag.Parse()
	// listening
	ln, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("wrong accept")
			continue
		}
		go handle(conn, *target)
	}

}

func handle(client net.Conn, target string) {
	defer client.Close()
	server, err := net.Dial("tcp", target)
	if err != nil {
		log.Printf("dial %s: %v", target, err)
		return
	}
	defer server.Close()
	done := make(chan struct{}, 2)
	go func() {
		io.Copy(server, client)
		done <- struct{}{}
	}()
	go func() {
		io.Copy(client, server)
		done <- struct{}{}
	}()

	<-done
}
