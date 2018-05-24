// Capturing stack traces
package main

import (
	"fmt"
	"runtime"
	"time"
)

func main() {
	go baz()
	foo()
}

func baz() {
	time.Sleep(1 * time.Second)
}

func foo() {
	bar()
}

func bar() {
	buf := make([]byte, 4096)
	l := runtime.Stack(buf, true /* false: only current goroutine */)
	fmt.Printf("Trace:\n%s\n", buf[:l])
}
