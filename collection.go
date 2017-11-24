package audioghost

import (
	"time"
	//"strings"
	//"log"
	"fmt"
)

type Collection struct {
	Name        string
	Path        string
	Audiobooks  Audiobooks
	Collections Collections
	Playtime    time.Duration
}

type Collections map[string]*Collection

//func (col *Collection) AddAudioBook(path string) error {
//	fmt.Println("Adding Book to:", col.Path)
//	col.Audiobooks.CreateAudioBook(path)
//	return nil
//}

func (col *Collection) AddAudioBook(book *Audiobook) error {
	fmt.Println("Adding Book to:", col.Name)
	col.Playtime += book.Playtime
	col.Audiobooks.AddAudioBook(book)
	return nil
}

func (col *Collection) RemoveAudioBook(book *Audiobook) error {
	fmt.Println("Removing %s from %s\n", book.Name, col.Name)
	col.Playtime -= book.Playtime
	col.Audiobooks.RemoveAudioBook(book)
	return nil
}

func NewCollection() Collection {
	return Collection{
		Audiobooks: make(map[string]*Audiobook),
		Collections: make(map[string]*Collection),
	}
}

func (col *Collection) String() {
	for _, f := range col.Audiobooks {
		fmt.Println("Name: ", f.Name);
		fmt.Println("Path: ", f.Path);
		fmt.Println("Duration:", f.Playtime)
		//fmt.Println("Files: ")
		//for _, mp3 := range f.Files {
		//	fmt.Print(mp3, ",")
		//}
		fmt.Println()
	}
}

func (col *Collection) AddToCollections(c *Collection) bool {
	col.Collections[c.Path] = c
	return true
}


func (col Collections) AddCollection(path string, name string) bool {
	for _, v := range col {
		if v.Path == path {
			return false
		}

	}
	fmt.Println("Added to Collections:", path)
	col[path] = &Collection{name, path, make(map[string]*Audiobook),make(map[string]*Collection),0}
	return true
}