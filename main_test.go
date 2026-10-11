package main

import (
	"io"
	"net"
	"testing"
	"time"
)

func startEcho(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0") // 0 = любой свободный порт
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return // листенер закрыт, тест закончился
			}
			go func() { defer c.Close(); io.Copy(c, c) }() // что пришло, то и отправляем назад
		}
	}()
	return ln.Addr().String() // например "127.0.0.1:53712"
}

func TestProxyEcho(t *testing.T) {
	target := startEcho(t)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })

	go serve(ln, target)
	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		t.Errorf("Failed dial connection")
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(2 * time.Second))

	want := "hello"
	if _, err := conn.Write([]byte(want)); err != nil {
		t.Fatal(err)
	}
	got := make([]byte, len(want))                    // буфер на 5 байт
	if _, err := io.ReadFull(conn, got); err != nil { // ждёт, пока придут все 5 байт
		t.Fatal(err)
	}
	if string(got) != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
