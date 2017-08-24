package main

/////////////////////////////////////////////////////////////////
//Code generated by chidley https://github.com/gnewton/chidley //
/////////////////////////////////////////////////////////////////

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	lib "github.com/gnewton/pstxml2sqlite/pstxml2sqlitestructs"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io"
	"log"
	"os"
	"strings"
)

const (
	JsonOut = iota
	XmlOut
	CountAll
)

var toJson bool = false
var toXml bool = false
var oneLevelDown bool = false
var countAll bool = false
var musage bool = false

var uniqueFlags = []*bool{
	&toJson,
	&toXml,
	&countAll}

var filename = "/home/gnewton/work/pst2json/foo.xml.gz"

func init() {
	flag.BoolVar(&toJson, "j", toJson, "Convert to JSON")
	flag.BoolVar(&toXml, "x", toXml, "Convert to XML")
	flag.BoolVar(&countAll, "c", countAll, "Count each instance of XML tags")
	flag.BoolVar(&oneLevelDown, "s", oneLevelDown, "Stream XML by using XML elements one down from the root tag. Good for huge XML files (see http://blog.davidsingleton.org/parsing-huge-xml-files-with-go/")
	flag.BoolVar(&musage, "h", musage, "Usage")
	flag.StringVar(&filename, "f", filename, "XML file or URL to read in")
}

var out int = -1

var counters map[string]*int

func main() {

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.CreateTable(&lib.Message{})

	flag.Parse()

	if musage {
		flag.Usage()
		return
	}

	numSetBools, outFlag := numberOfBoolsSet(uniqueFlags)
	if numSetBools == 0 {
		flag.Usage()
		return
	}

	if numSetBools != 1 {
		flag.Usage()
		log.Fatal("Only one of ", uniqueFlags, " can be set at once")
	}

	reader, xmlFile, err := genericReader(filename)
	if err != nil {
		log.Fatal(err)
		return
	}

	counter := 0
	tx := db.Begin()

	decoder := xml.NewDecoder(reader)
	counters = make(map[string]*int)
	for {

		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch se := token.(type) {
		case xml.StartElement:
			handleFeed(se, decoder, outFlag, tx)
		}
		counter = counter + 1
		if counter == 5000 {
			tx.Commit()
			counter = 0
			tx = db.Begin()
		}
	}
	if counter > 0 {
		tx.Commit()
	}
	if xmlFile != nil {
		defer xmlFile.Close()
	}
	if countAll {
		for k, v := range counters {
			fmt.Println(*v, k)
		}
	}
}

func handleFeed(se xml.StartElement, decoder *xml.Decoder, outFlag *bool, db *gorm.DB) {
	if se.Name.Local == "message" && se.Name.Space == "" {
		var item lib.Message
		decoder.DecodeElement(&item, &se)
		switch outFlag {
		case &toJson:
			//writeJson(item)
		case &toXml:
			//	writeXml(item)
		}
		fmt.Println(item)
		fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
		//db.Create(item)
		db.Save(&item)
	}

	if se.Name.Local == "meta" && se.Name.Space == "" {
		var item lib.Meta
		decoder.DecodeElement(&item, &se)
		switch outFlag {
		case &toJson:
			//writeJson(item)
		case &toXml:
			//writeXml(item)
		}
	}
}

func makeKey(space string, local string) string {
	if space == "" {
		space = "_"
	}
	return space + ":" + local
}

func incrementCounter(space string, local string) {
	key := makeKey(space, local)

	counter, ok := counters[key]
	if !ok {
		n := 1
		counters[key] = &n
	} else {
		newv := *counter + 1
		counters[key] = &newv
	}
}

func writeJson(item interface{}) {
	b, err := json.MarshalIndent(item, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func writeXml(item interface{}) {
	output, err := xml.MarshalIndent(item, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	os.Stdout.Write(output)
}

func genericReader(filename string) (io.Reader, *os.File, error) {
	if filename == "" {
		return bufio.NewReader(os.Stdin), nil, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	if strings.HasSuffix(filename, "bz2") {
		return bufio.NewReader(bzip2.NewReader(bufio.NewReader(file))), file, err
	}

	if strings.HasSuffix(filename, "gz") {
		reader, err := gzip.NewReader(bufio.NewReader(file))
		if err != nil {
			return nil, nil, err
		}
		return bufio.NewReader(reader), file, err
	}
	return bufio.NewReader(file), file, err
}

func numberOfBoolsSet(a []*bool) (int, *bool) {
	var setBool *bool
	counter := 0
	for i := 0; i < len(a); i++ {
		if *a[i] {
			counter += 1
			setBool = a[i]
		}
	}
	return counter, setBool
}
