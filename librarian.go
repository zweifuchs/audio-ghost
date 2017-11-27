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

//var db *sql.DB

type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}
func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Librarian struct {
	conf        *config.Config
	audiobooks  Audiobooks
	collections Collections
}

func sortCollections(collections *Collections) {
	var cpaths []string
	for k, _ := range *collections {
		cpaths = append(cpaths, k)
	}
	fmt.Println(cpaths)
	sort.Sort(ByLength(cpaths))
	fmt.Println(cpaths)

	for _, c := range cpaths {
		for _, cc := range cpaths {
			if strings.HasPrefix(c, cc) && c != cc {
				for _, tst := range *collections {
					if _, ok := tst.Collections[cc]; ok {
						//sortCollections(&tst.Collections)
						delete(tst.Collections, c)
					}
				}
				(*collections)[cc].AddToCollections((*collections)[c])
			}
		}
	}
}

func sortBooksAndCollections(books *Audiobooks, collections *Collections) {
	fmt.Println("Sorting....")

	var cpaths []string
	for k, _ := range *collections {
		cpaths = append(cpaths, k)
	}
	fmt.Println(cpaths)
	sort.Sort(ByLength(cpaths))
	fmt.Println(cpaths)

	// sort audiobooks first
	for _, c := range cpaths {
		for _, b := range *books {
			fmt.Printf("\"%s\" ::: \"%s\"\n", b.Path, c)
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

	// sort collections
	sortCollections(collections)
}

func (l *Librarian) Init(c *config.Config) {
	fmt.Println("Hello I am you local Librarian.")
	l.conf = c
	directory := l.conf.RootDirecotry()

	l.audiobooks = make(Audiobooks)
	l.collections = make(Collections)

	l.collections.AddCollection(l.conf.RootDirecotry(), "root")
	err := ScanDir(l, directory)
	if err != nil {
		log.Fatal(err)
	}

	sortBooksAndCollections(&l.audiobooks, &l.collections)

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
