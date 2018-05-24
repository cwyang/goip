// Verifying interfaces with canary tests
// When you're writing interfaces or implemetations of interfaces...

package main

import (
	"fmt"
	"io"
	"testing"
)

type MyWriter struct {
}

func (m *MyWriter) Write([]byte) error {
	return nil
}

func main() {
	m := map[string]interface{}{
		"w": &MyWriter{},
	}
	doSomething(m)
}

func doSomething(m map[string]interface{}) {
	w, ok := m["w"].(io.Writer)
	fmt.Println(w, ok)
}

func TestWriter(t *testing.T) {
	var _ io.Writer = &MyWriter{}
}
