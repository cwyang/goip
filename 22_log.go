// Logging to an arbitrary writer
package main

import (
	"log"
	"os"
)

func main() {
	//logfile, _ := os.Create("./log.txt")
	logfile, _ := os.OpenFile("./log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer logfile.Close()

	logger := log.New(logfile, "example ", log.LstdFlags|log.Lshortfile)
	logger.Println("Hello")
	logger.Fatalln("Second Hello")
	logger.Println("Final Hello")
}
