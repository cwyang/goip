// Goroutine with panic and recover
package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

func main() {
	listen()
}

func listen() {
	listener, e := net.Listen("tcp", ":1026")
	if e != nil {
		fmt.Println("listener error")
		return
	}
	for {
		conn, e := listener.Accept()
		if e != nil {
			fmt.Println("accept error")
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("Fatal error: %s\n", e)
		}
		conn.Close()
	}()
	reader := bufio.NewReader(conn)

	data, e := reader.ReadBytes('\n')
	if e != nil {
		fmt.Println("read err")
		conn.Close()
	}
	response(data, conn)
}

func response(data []byte, conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	conn.Write(data)
	panic(errors.New("panic test"))
}
