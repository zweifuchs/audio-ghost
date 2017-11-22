package audioghost

import (
	"fmt"
	"path/filepath"
	"log"
	"strings"
	"os"
	"time"
	"github.com/tcolgate/mp3"
	// dd "github.com/zweifuchs/audioghost/debug"
)

type Audiobook struct {
	Name  string
	Path  string
	Files []string
	Playtime time.Duration
}

type Audiobooks map[string]*Audiobook

func (book *Audiobook) AddFiles() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}
		if strings.HasSuffix(info.Name(), ".mp3") {
			length, _ := calcmp3length(path)
			fmt.Printf("Adding %s to %s playtime %s \n", path, book.Name, length)
			book.Playtime += length
			book.Files = append(book.Files, path)
		}
		//fmt.Println(book)
		return nil
	}
}

func calcmp3length(path string) (time.Duration, error) {
	skipped := 0
	var length time.Duration = 0

	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	d := mp3.NewDecoder(file)
	var f mp3.Frame
	for {

		if err := d.Decode(&f, &skipped); err != nil {
			//fmt.Println(err)
			break
		}
		length += f.Duration()
		//fmt.Println(&f)
	}

	return length, nil

}



func (audiobooks Audiobooks) isDuplicate (path string) bool {
	for k, _ := range audiobooks {
		if strings.HasPrefix(path, k) {
			return true
		}
	}
	return false
}

func (audiobooks Audiobooks) CreateAudioBook (path string) error {
	bookName := getBookName(path)

	if _, ok := audiobooks[path]; ok {
		fmt.Printf("%v already in list\n", bookName)
		return nil
	} else {
		// check if path is already used
		if  audiobooks.isDuplicate(path) {
			fmt.Printf("%v already in list\n", bookName)
			return nil
		}
		fmt.Printf("Add %v to list\n", bookName)
		tmpbook := Audiobook{
			Name: bookName,
			Path: path,
		}

		err := filepath.Walk(path, tmpbook.AddFiles())
		if err != nil {
			log.Fatal(err)
		}

		audiobooks[path] = &tmpbook
		//fmt.Println("New book...", tmpbook)
	}
	return nil
}

