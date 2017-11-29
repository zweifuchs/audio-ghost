package audioghost

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"database/sql"
	_ "golang.org/x/tools/cmd/guru/testdata/src/lib"
)

func getBookName(path string) string {
	//currentPath := filepath.Dir(path)
	currentBook := filepath.ToSlash(path)
	return currentBook[strings.LastIndex(currentBook, "/")+1:]
}

func getLastPathDir(path string) string {
	return getBookName(path)
}

func ScanDir(librarian *Librarian ,llbooks *Audiobooks, llcollections *Collections, path string) error {

	books_ := make(Audiobooks)
	collections_ := make(Collections)

	books := &books_
	collections := &collections_

	err := filepath.Walk(path, walkDir(books, collections))

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
	db.Exec("DROP TABLE files;");
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS audiobooks
		 (
		    ID INT NOT NULL AUTO_INCREMENT,
		    PRIMARY KEY(ID),
		    Name VARCHAR(255),
		    Path VARCHAR(2047),
		    Files TEXT,
		    Playtime BIGINT,
		    Description TEXT
		 ) ENGINE=INNODB;

	`)
	checkErr(err)
	_ , err = db.Exec(`
		CREATE TABLE IF NOT EXISTS collections
		 (
		    ID INT NOT NULL AUTO_INCREMENT,
		    PRIMARY KEY(ID),
		    Name VARCHAR(255),
		    Path VARCHAR(2047),
		    Playtime BIGINT
		 )  ENGINE=INNODB;
;
	`)
	checkErr(err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS files
		 (
		    ID INT NOT NULL AUTO_INCREMENT,
		    PRIMARY KEY(ID),
		    Path VARCHAR(2047),
		    Playtime BIGINT,
		    Audiobook_id INT NOT NULL,
		    INDEX audiobook_ind1 (Audiobook_id),
		    FOREIGN KEY (Audiobook_id)
        	    	REFERENCES audiobooks(ID)
		 ) ENGINE=INNODB;
	`)
	checkErr(err)
	/*
	INSERT AUDIOBOOKS
	*/

	//for _book := range l.audiobooks {
	//	fmt.Printf("%t, %s \n",book, book)
	//}

	for _, book := range *books {
		_, err := db.Exec("INSERT INTO audiobooks (Name,Path,Files,Playtime,Description) VALUES (?,?,?,?,?)",
			book.Name,
			book.Path,
			strings.Join(book.Files, string(librarian.conf.PathListSeperator)),
			book.Playtime,
			book.Description,
		)
		checkErr(err)
		var audiobook_id int
		row,err := db.Query("SELECT ID FROM audiobooks WHERE Name=?",book.Name)
		checkErr(err)

		defer row.Close()
		row.Next()
		row.Scan(&audiobook_id)

		fmt.Println(audiobook_id)
		for _, file := range book.Files {
			_, err = db.Exec("INSERT INTO files (Path,Playtime,Audiobook_id) VALUES (?,?,?)",
				file,
				0,
				audiobook_id,
			)
			checkErr(err)
		}

		row.Close()
	}

	for _, collection := range *collections {
		_, err := db.Exec("INSERT INTO collections (Name,Path,Playtime) VALUES (?,?,?)",
			collection.Name,
			collection.Path,
			collection.Playtime,
		)
		checkErr(err)
	}

	return err
}

func walkDir(books *Audiobooks, collections *Collections) filepath.WalkFunc {

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
				log.Print(err)
				return err
			}

			for _, v := range files {
				if v.Name() == "ab_root" || strings.HasSuffix(v.Name(), "mp3") {
					if v.Name() == "ab_root" {
						fmt.Println("Found a root file:", path)
					}
					books.CreateAudioBook(path)
					return nil
				}
				if v.Name() == "collection" {
					fmt.Println("Found a New Collection:", path)
					collections.CreateAndAddCollection(path, getLastPathDir(path))
					return nil
				}
			}

		}
		return nil
	}
}

func getBooks(librarian *Librarian, audiobooks *Audiobooks) error{
	db, err := sql.Open("mysql", "audioghost:123456@/audioghost")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM audiobooks")
	checkErr(err)

	for rows.Next() {
		var r = new(Audiobook)
		err = rows.Scan(&r.Id, &r.Name, &r.Path,&r.FilesAsText,&r.Playtime,&r.Description)
		audiobooks.AddAudioBook(r)
	}
	return err
}

func getCollections(librarian *Librarian, collections *Collections) error {
	db, err := sql.Open("mysql", "audioghost:123456@/audioghost")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM collections")
	checkErr(err)

	for rows.Next() {
		var r = NewCollection()
		err = rows.Scan(&r.Id, &r.Name, &r.Path, &r.Playtime, )
		collections.AddCollection(&r)
	}
	return err
}