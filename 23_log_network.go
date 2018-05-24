// Network logging
// use netcat: nc -lk 8888 (tcp) or nc -lku 8888 (udp)
package main

import (
	"log"
	"net"
	"time"
)

func main() {
	//	conn, e := net.Dial("tcp", "localhost:8888")
	timeout := 3 * time.Second
	conn, e := net.DialTimeout("udp", "localhost:8888", timeout)
	if e != nil {
		panic("connect")
	}
	defer conn.Close()

	f := log.Ldate | log.Lshortfile
	logger := log.New(conn, "test ", f)
	logger.Println("regular mesg")
	logger.Panicln("panic mesg")
}
