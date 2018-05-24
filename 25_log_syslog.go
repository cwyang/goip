// Logging to the syslog
package main

import (
	"fmt"
	//	"log"
	"log/syslog"
)

func main() {
	// prio := syslog.LOG_LOCAL3 | syslog.LOG_NOTICE
	// flags := log.Ldate | log.Lshortfile
	// logger, e := syslog.NewLogger(prio, flags)
	logger, e := syslog.New(syslog.LOG_LOCAL3, "test")
	if e != nil {
		fmt.Printf("syslog err: %s", e)
		return
	}
	defer logger.Close()

	logger.Debug("debug")
	logger.Notice("notice")
	logger.Warning("warn")
	logger.Alert("alert")
	//	logger.Println("test log")
}
