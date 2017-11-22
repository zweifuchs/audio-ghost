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

func ScanDir(path string, books Audiobooks) error {
	err := filepath.Walk(path, walkDir(books))
	return err
}

func walkDir(audiobooks Audiobooks) filepath.WalkFunc {

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
				log.Print(err); return err; }
			for _, v := range files {
				if v.Name() == "ab_root" || strings.HasSuffix(v.Name(), "mp3") {
					fmt.Println("Found a root file:", path)
					audiobooks.CreateAudioBook(path)
					return nil
				}
			}

		}
		return nil
	}
}