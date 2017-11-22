package main

import (
	"fmt"
	_ "io/ioutil"
	"log"
	"flag"
	_ "os"
	_ "strings"
	ag "github.com/zweifuchs/audioghost"
	dd "github.com/zweifuchs/audioghost/debug"
	"time"
)

func main() {
	defer dd.TimeTracker(time.Now(),"main")
	var directory string
	flag.StringVar(&directory, "dir", "./", "directory to scan")
	flag.Parse()
	fmt.Println("Scanning directory", directory)

	books := make(ag.Audiobooks)
	collections := make(ag.Collections)
	bigCollection := ag.Collection{
		Name: "all",
		Path: directory,
		Audiobooks: &books,
		Collections: &collections,
		Playtime: 0,
	}

	err := ag.ScanDir(bigCollection, directory)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n\nResults:")
	for _, f := range *bigCollection.Audiobooks {
		fmt.Println("Name: ",f.Name);
		fmt.Println("Path: ", f.Path);
		fmt.Println("Duration:", f.Playtime)
		//fmt.Println("Files: ")
		//for _, mp3 := range f.Files {
		//	fmt.Print(mp3, ",")
		//}
		fmt.Println()
	}

}
