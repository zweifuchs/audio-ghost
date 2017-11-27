package audioghost

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zweifuchs/audioghost/lib/config"
	"log"
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

	l.collections.AddCollection(l.conf.RootDirecotry(), "root")


	// Force Rescan?
	if l.conf.Rescan {
		err := ScanDir(l, directory)
		if err != nil {
			log.Fatal(err)
		}
	}

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

	/*
	 DB STUFF
	*/
	fmt.Println("Saving to DB:")
	db, err := sql.Open("mysql", "audioghost:123456@/audioghost")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		// do something here
	}

	/*
	CREAT TABLE
	*/
	db.Exec("DROP TABLE audiobooks;");
	db.Exec("DROP TABLE collections;");
	res, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS audiobooks
		 (
		    ID INT NOT NULL AUTO_INCREMENT,
		    PRIMARY KEY(ID),
		    Name VARCHAR(255),
		    Path VARCHAR(2047),
		    Files TEXT,
		    Playtime BIGINT,
		    Description TEXT
		 );
	`)
	res, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS collections
		 (
		    ID INT NOT NULL AUTO_INCREMENT,
		    PRIMARY KEY(ID),
		    Name VARCHAR(255),
		    Path VARCHAR(2047),
		    Playtime BIGINT
		 );
	`)
	if err != nil {
		panic(err.Error())
	}

	/*
	INSERT AUDIOBOOKS
	*/

	//for _book := range l.audiobooks {
	//	fmt.Printf("%t, %s \n",book, book)
	//}

	for _,book := range l.audiobooks {
		_, err := db.Exec("INSERT INTO audiobooks (Name,Path,Files,Playtime,Description) VALUES (?,?,?,?,?)",
			book.Name,
			book.Path,
			strings.Join(book.Files,","),
			book.Playtime,
			book.Description,
		)
		checkErr(err)
	}

	for _, collection := range l.collections {
		_, err := db.Exec("INSERT INTO collections (Name,Path,Playtime) VALUES (?,?,?)",
			collection.Name,
			collection.Path,
			collection.Playtime,
		)
		checkErr(err)
	}

	//	//defer stmt.Close()
	//
	//
	//
	//	//stmt.Close()
	//}

	fmt.Println("Db result:", res)

}
