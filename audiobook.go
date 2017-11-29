package audioghost

import (
	"fmt"
	"github.com/tcolgate/mp3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	// dd "github.com/zweifuchs/audioghost/debug"
)

type Audiobook struct {
	Id          int
	Name        string
	Path        string
	FilesAsText string
	Files       []MediaFile
	Playtime    time.Duration
	Description string
}

type Audiobooks map[string]*Audiobook

func (book *Audiobook) AddFiles() filepath.WalkFunc {
	var trcknr int = 0
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}
		if strings.HasSuffix(info.Name(), ".mp3") {
			length, _ := calcmp3length(path)
			fmt.Printf("Adding %s to %s playtime %s \n", path, book.Name, length)
			book.Playtime += length
			book.Files = append(book.Files, MediaFile{Name: book.Name, Path: path, AudioBook: book ,Track: trcknr,Playtime: time.Duration(length)})
		}
		//fmt.Println(book)
		trcknr += 1
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

func (audiobooks Audiobooks) isDuplicate(path string) bool {
	for k, _ := range audiobooks {
		if strings.HasPrefix(path, k) {
			return true
		}
	}
	return false
}

func (audiobooks Audiobooks) CreateAudioBook(path string) error {
	bookName := getBookName(path)

	if _, ok := audiobooks[path]; ok {
		fmt.Printf("%v already in list\n", bookName)
		return nil
	} else {
		// check if path is already used
		if audiobooks.isDuplicate(path) {
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

func (audiobooks Audiobooks) AddAudioBook(book *Audiobook) error {
	audiobooks[book.Path] = book
	return nil
}
func (audiobooks Audiobooks) RemoveAudioBook(book *Audiobook) error {
	delete(audiobooks, book.Path)
	return nil
}
