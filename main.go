package main

/////////////////////////////////////////////////////////////////
//Initial code generated by chidley https://github.com/gnewton/chidley //
/////////////////////////////////////////////////////////////////

import (
	"database/sql"
	"encoding/xml"
	lib "github.com/gnewton/pstxml2sqlite/pstxml2sqlitestructs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"sync"
)

const (
	JsonOut = iota
	XmlOut
	CountAll
)

//var filenames = []string{"/home/gnewton/aafc_email_pst/all.xml.bz2"}

//var filenames = []string{"/home/gnewton/work/pst2json/all.xml.bz2"}

//var filenames = []string{"/home/gnewton/work/pst2json/archive_2018_June15.pst.xml.bz2"}

var filenames = []string{"/home/gnewton/work/pst2json/archive.xml.bz2"}

//var filenames = []string{"/home/gnewton/work/pst2json/m"}

//
var counter = 0
var attachCounter int64 = 0
var recipCounter int64 = 0
var countAll int64 = 0
var countAll2 int64 = 0

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// ADD check to see if file exists
	// use gorm to create DB
	db, err := gorm.Open("sqlite3", "a.db")
	if err != nil {
		panic("failed to connect database")
	}
	//defer db.Close()

	db.CreateTable(&lib.Message{})
	db.CreateTable(&lib.Recipient{})
	db.CreateTable(&lib.Attachment{})
	db.CreateTable(&lib.Filesource{})
	db.Close()

	db2, err := sql.Open("sqlite3", "file:a.db")
	if err != nil {
		log.Fatal(err)
	}

	err = setPragmas(db2)
	if err != nil {
		log.Fatal(err)
	}

	tx2, err := db2.Begin()
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan []*lib.Message, 200)

	messageStmt, err := newStatement(tx2, messageSql)
	if err != nil {
		log.Fatal(err)
	}

	attachStmt, err := newStatement(tx2, attachmentSql)
	if err != nil {
		log.Fatal(err)
	}

	recipStmt, err := newStatement(tx2, recipSql)
	if err != nil {
		log.Fatal(err)
	}

	var mux sync.Mutex

	saver := &MessageSaver{
		tx2:         tx2,
		db2:         db2,
		c:           c,
		messageStmt: messageStmt,
		attachStmt:  attachStmt,
		recipStmt:   recipStmt,
		mux:         &mux,
	}

	var wg sync.WaitGroup

	num := 2

	wg.Add(num)
	for i := 0; i < num; i++ {
		go saver.run(&wg)
	}

	dups = make(map[string]bool, 0)

	for i := 0; i < len(filenames); i++ {
		filename := filenames[i]
		log.Println("Opening file: " + filename)

		reader, _, err := genericReader(filename)
		if err != nil {
			log.Fatal(err)
			return
		}
		decoder := xml.NewDecoder(reader)

		for {
			token, _ := decoder.Token()
			if token == nil {
				break
			}
			switch se := token.(type) {
			case xml.StartElement:
				handleFeed(saver, se, decoder, c)
			}
		}
	}
	endCommit := false
	if chunk != nil {
		c <- chunk
		endCommit = true
	}
	close(c)
	wg.Wait()

	if counter > 0 || endCommit {
		log.Println("Final commit")
		err := saver.tx2.Commit()
		if err != nil {
			log.Fatal(err)
		}

	}

	log.Println("Total messages")
	log.Println(countAll)
	log.Println(countAll2)

}
