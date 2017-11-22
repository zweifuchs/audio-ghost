package main

import (
	"fmt"
	_ "io/ioutil"
	"log"
	"flag"
	_ "os"
	"path/filepath"
	_ "strings"
	ag "github.com/zweifuchs/audioghost"
	dd "github.com/zweifuchs/audioghost/debug"
	"time"
)

func main() {
	defer dd.TimeTracker(time.Now(),"main")
	var directory string
	books := make(ag.Audiobooks)

	flag.StringVar(&directory, "dir", "./", "directory to scan")
	flag.Parse()

	fmt.Println("Scanning directory", directory)

	err := filepath.Walk(directory, ag.Scandir(books))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n\nResults:")
	for _, f := range books {
		fmt.Println("Name: ",f.Name);
		fmt.Println("Path: ", f.Path);
		fmt.Println("Files: ")
		for _, mp3 := range f.Files {
			fmt.Print(mp3, ",")
		}
		fmt.Println()
	}

}
