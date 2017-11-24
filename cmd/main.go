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

	var c config.Config
	c.ReadCmd()
	directory := c.RootDirecotry()
	librarian := new(ag.Librarian)
	fmt.Println("Scanning directory", directory)
	librarian.Init(&c)

}
