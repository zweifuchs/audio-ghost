package audioghost

import (
	"time"
	"strings"
	"log"
	"fmt"
)

type Collection struct {
	Name        string
	Path        string
	Audiobooks  *Audiobooks
	Collections *Collections
	Playtime    time.Duration
}


type Collections map[string]*Collection

func (col *Collection) AddAudioBook(audioBooks Audiobooks, path string) error {

	return nil
}

func (col *Collection) IsDuplicate(path string) bool {
	for k, _ := range *col.Audiobooks {
		if strings.HasPrefix(path, k) {
			return true
		}
	}

	for _,v := range *col.Collections {
		if v.IsDuplicate(path) {
			return true
		}

	}

	return false
}

func (col Collections) AddCollection(path string) {

	for _, v := range col {

		if v.Path == path {
			return
		}

	}

	if _, ok := col[path]; ok {
		fmt.Printf("%v already in list\n", path)
		return
	}
	books := make(Audiobooks)
	collections := make(Collections)
	result := Collection{
		Name: path,
		Path: path,
		Audiobooks: &books,
		Collections: &collections,
		Playtime: 0,
	}
	col[path] = &result
	err := ScanDir(result, path)
	if err != nil {
		log.Fatal(err)
	}
}