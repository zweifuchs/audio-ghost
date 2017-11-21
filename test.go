package main

import (
	"fmt"
	_ "io/ioutil"
	"log"
	"flag"
	"os"
	"path/filepath"
	"strings"
)

type Audiobook struct {
	Name  string
	Path  string
	Files []string
}

type Audiobooks map[string]*Audiobook

type Collection struct {
	Name string
	Path string
	Audiobooks Audiobooks
	Collections Collections
}

type Collections map[string]*Collection

func consume(interface{}) int {
	return 0
}

func scandir(audiobooks Audiobooks) filepath.WalkFunc {

	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}
		if strings.HasSuffix(info.Name(), ".mp3") {
			fmt.Printf("Found Mp3 %v in Folder %v \n", info.Name(), path)
			//currentpath := path[0:strings.LastIndex(path, info.Name())]
			currentpath := filepath.Dir(path)
			currentbook := filepath.ToSlash(currentpath)
			currentbook = currentbook[strings.LastIndex(currentbook, "/") + 1:]

			// Check if we already have that book
			if val, ok := audiobooks[currentpath]; ok {
				val.Files = append(val.Files, info.Name())
			} else {
				tmpbook := Audiobook{
					Name: currentbook,
					Path: currentpath,

				}
				tmpbook.Files = append(tmpbook.Files, info.Name())

				audiobooks[currentpath] = &tmpbook
				fmt.Println("New book...", tmpbook)
			}
		}
		return nil
	}
}

func main() {

	var directory string
	books := make(Audiobooks)

	flag.StringVar(&directory, "dir", "./", "directory to scan")
	flag.Parse()

	fmt.Println("Scanning directory", directory)

	err := filepath.Walk(directory, scandir(books))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Results:")
	for _, f := range books {
		fmt.Println(f.Name);
		fmt.Println(f.Path);
		for _, mp3 := range f.Files {
			fmt.Print(mp3, ",")
		}
		fmt.Println()
	}

}
