package audioghost

import (
	"github.com/zweifuchs/audioghost/lib/config"
	"fmt"
	"log"
	"strings"
)

type Librarian struct {
	conf *config.Config
	audiobooks Audiobooks
	collections Collections
}

func sortBooksAndCollections(books *Audiobooks, collections *Collections) {
	fmt.Println("Sorting....")
	for _, c := range *collections {
		for _, b := range *books {
			fmt.Printf("\"%s\" ::: \"%s\"\n",b.Path, c.Path)
			if strings.HasPrefix(b.Path,c.Path) {
				c.AddAudioBook(b)
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

	err := ScanDir(l, directory)
	if err != nil {
		log.Fatal(err)
	}

	sortBooksAndCollections(&l.audiobooks, &l.collections)

	fmt.Println("\n\nBooks:")
	for _,v := range l.audiobooks {
		fmt.Printf("%s,%s,\n",v.Name, v.Playtime)
	}

	fmt.Println("\nCollections:")
	for _, v := range l.collections {
		fmt.Printf("%s mit \n", v.Name)
		fmt.Println("===========")
		for _,b := range v.Audiobooks {
			fmt.Println(b.Name)
		}
	}
}