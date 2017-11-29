package audioghost

import (
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zweifuchs/audioghost/lib/config"
	//"log"
	"sort"
	"strings"
)

type Librarian struct {
	conf        *config.Config
	audiobooks  Audiobooks
	collections Collections
}

// Sorting collections
// this takes a map of collections [PATH]Collection and puts the nexted one to the right place
func sortCollections(collections *Collections) {
	var cpaths []string
	for k, _ := range *collections {
		cpaths = append(cpaths, k)
	}

	// sorting. ByLength is defined in helper.go
	sort.Sort(ByLength(cpaths))
	for _, c := range cpaths {
		for _, cc := range cpaths {
			if strings.HasPrefix(c, cc) && c != cc {
				for _, tst := range *collections {
					if _, ok := tst.Collections[cc]; ok {
						delete(tst.Collections, c)
					}
				}
				(*collections)[cc].AddToCollections((*collections)[c])
			}
		}
	}
}

func sortBooks(books *Audiobooks, collections *Collections) {
	var cpaths []string
	for k, _ := range *collections {
		cpaths = append(cpaths, k)
	}

	// sorting. ByLength is defined in helper.go
	sort.Sort(ByLength(cpaths))

	for _, c := range cpaths {
		for _, b := range *books {
			if strings.HasPrefix(b.Path, c) {
				for _, tst := range *collections {
					if _, ok := tst.Audiobooks[b.Path]; ok {
						tst.RemoveAudioBook(b)
					}
				}
				(*collections)[c].AddAudioBook(b)
			}
		}
	}
}

func (l *Librarian) Init(c *config.Config) {
	fmt.Println("Hello I am you local Librarian.")
	l.conf = c
	directory := l.conf.RootDirecotry()

	l.audiobooks = make(Audiobooks)
	l.collections = make(Collections)

	l.collections.CreateAndAddCollection(l.conf.RootDirecotry(), "root")


	// Force Rescan?
	if l.conf.Rescan {
		err := ScanDir(l, &l.audiobooks, &l.collections, directory)
		checkErr(err)
	}

	err := getBooks(l, &l.audiobooks)
	checkErr(err)
	err = getCollections(l, &l.collections)
	checkErr(err)


	// Get DB Data

	sortBooks(&l.audiobooks, &l.collections)
	sortCollections(&l.collections)

	fmt.Println("\n\nBooks:")
	for _, v := range l.audiobooks {
		fmt.Printf("%s,%s,\n", v.Name, v.Playtime)
	}

	fmt.Println("\nCollections:")
	for _, v := range l.collections {
		fmt.Printf("\n\n%s mit \n", v.Name)
		fmt.Println("===========")
		for _, b := range v.Collections {
			fmt.Println("Collection:", b.Name)
		}
		for _, b := range v.Audiobooks {
			fmt.Println("Audiobook:", b.Name)
		}
	}



}

func (l *Librarian) GetFile() {

}