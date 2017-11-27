package main

import (
	"fmt"
	ag "github.com/zweifuchs/audioghost"
	dd "github.com/zweifuchs/audioghost/lib/debug"
	"github.com/zweifuchs/audioghost/lib/config"
	"time"
)

func main() {
	defer dd.TimeTracker(time.Now(),"main")

	librarian := new(ag.Librarian)
	librarian.Init(config.GetConfig())

	fmt.Println("Thank you for using our service ...")
}
