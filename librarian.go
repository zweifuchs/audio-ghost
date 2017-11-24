package audioghost

import (
	"github.com/zweifuchs/audioghost/lib/config"
	"fmt"
	"log"
	"strings"
)

type allAudioBooks struct {
	audiobooks Audiobooks
}

type allCollections struct {
	collections Collections
}

type Librarian struct {
	conf *config.Config
	allAudioBooks allAudioBooks
	allCollections allCollections
}

func sortBooksAndCollections(books *Audiobooks, collections *Collections) {
	for _, c := range *collections {
		for _, b := range *books {
			if strings.HasPrefix(b.Path,c.Path) {

			}
		}
	}
}

func (l *Librarian) Init(c *config.Config) {
	fmt.Println("Hello I am you local Librarian.")
	l.conf = c
	directory := l.conf.RootDirecotry()

	l.allAudioBooks.audiobooks = make(Audiobooks)
	l.allCollections.collections = make(Collections)

	err := ScanDir(l, directory)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n\nBookss:")
	for _,v := range l.allAudioBooks.audiobooks {
		fmt.Printf("%s mit %s\n",v.Name, v.Playtime)
	}

	fmt.Println("\nCollections:")
	for _, v := range l.allCollections.collections {
		fmt.Printf("%s mit \n", v.Name)
		fmt.Println("===========")
		for _,b := range v.Audiobooks {
			fmt.Println(b.Name)
		}
	}
}