package audioghost

import (
	"fmt"
	"path/filepath"
	"log"
	"os"
	"io/ioutil"
	"strings"
)

func getBookName(path string) string {
	//currentPath := filepath.Dir(path)
	currentBook := filepath.ToSlash(path)
	return currentBook[strings.LastIndex(currentBook, "/") + 1:]
}

func addFiles(book *Audiobook) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}
		if strings.HasSuffix(info.Name(), ".mp3") {
			fmt.Printf("Adding %s to %s \n", path, book.Name)
			(*book).Files = append((*book).Files, path)
		}
		//fmt.Println(book)
		return nil
	}
}

func isNestedPath(audiobooks Audiobooks, path string) bool {
	for k, _ := range audiobooks {
		if strings.HasPrefix(path, k) {
			return true
		}
	}

	return false
}

func CreateAudioBook(audiobooks Audiobooks, path string) error {
	bookName := getBookName(path)

	if _, ok := audiobooks[path]; ok {
		fmt.Printf("%v already in list\n", bookName)
		return nil
	} else {
		// check if path is already used
		if isNestedPath(audiobooks, path) {
			fmt.Printf("%v already in list\n", bookName)
			return nil
		}
		fmt.Printf("Add %v to list\n", bookName)
		tmpbook := Audiobook{
			Name: bookName,
			Path: path,
		}

		err := filepath.Walk(path, addFiles(&tmpbook))
		if err != nil {
			log.Fatal(err)
		}

		audiobooks[path] = &tmpbook
		//fmt.Println("New book...", tmpbook)
	}
	return nil
}

func Scandir(audiobooks Audiobooks) filepath.WalkFunc {

	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}

		// Look inside if this might be a root folder
		//fmt.Println("Checking:",path)
		if info.IsDir() {
			files, err := ioutil.ReadDir(path)
			if err != nil {
				log.Print(err); return nil; }
			for _, v := range files {
				if v.Name() == "ab_root" || strings.HasSuffix(v.Name(), "mp3") {
					fmt.Println("Found a root file:", path)
					CreateAudioBook(audiobooks, path)
					return nil
				}
			}

		}
		return nil
	}
}