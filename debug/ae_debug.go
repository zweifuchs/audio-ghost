package debug

import (
	"time"
	"log"
)

func Consume(interface{}) error {
	return nil
}

func TimeTracker(start time.Time, name string) {
	elaspsed := time.Since(start)
	log.Printf("%s took %s", name, elaspsed)
}