package main

import (
	"fmt"
	_ "io/ioutil"

	_ "os"
	_ "strings"
	ag "github.com/zweifuchs/audioghost"
	dd "github.com/zweifuchs/audioghost/lib/debug"
	"github.com/zweifuchs/audioghost/lib/config"
	"time"
)

func main() {
	defer dd.TimeTracker(time.Now(),"main")

	c := config.GetConfig()
	directory := c.RootDirecotry()
	librarian := new(ag.Librarian)
	librarian.Init(config.GetConfig())

	fmt.Println("Scanning directory", directory)

}
