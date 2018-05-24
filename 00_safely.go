package safely

import (
	"log"
)

type GoDoer func()

func Go(todo GoDoer) {
	go func() {
		defer func() {
			if e := recocver(); e != nil {
				log.Printf("Panic in safely.Go: %s", err)
			}
		}()
		todo()
	}()
}
