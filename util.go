package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"database/sql"
	"encoding/base64"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func decodeBase64(s string) ([]byte, error) {
	val, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		log.Println(err)
	}
	return val, err
}

func string2bool(s string) bool {
	if s != "" && (s == "true" || s == "TRUE") {
		return true
	}
	return false

}

const baseDateLayout = "02-01-2006 06:04:05"
const altDateLayout = "01-02-2006 06:04:05"

const missingDateString = "01-01-1969 00:00:01"

var missingDate *time.Time

func string2date(str string) (*time.Time, error) {

	if str == "" {
		if missingDate == nil {
			t, err := time.Parse(baseDateLayout, missingDateString)
			if err != nil {
				log.Println(err)
			}
			missingDate = &t
		}
		return missingDate, nil
	}
	t, err := time.Parse(baseDateLayout, str)

	if err != nil {
		log.Println(err)
		log.Println(str)
		t, err = time.Parse(altDateLayout, str)
		if err != nil {
			log.Println(err)
			log.Println(str)
		} else {
			log.Println("Parsed ok, witd DD-MM-YY")
			log.Println(str)
		}
	}
	return &t, err
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

var pragmas = []string{"PRAGMA synchronous=off", "PRAGMA journal_mode=off", "PRAGMA cache_size=20000", "PRAGMA locking_mode=exclusive", "PRAGMA temp_store=memory"}

func setPragmas(db *sql.DB) error {

	for i := 0; i < len(pragmas); i++ {
		log.Println(pragmas[i])
		_, err := db.Exec(pragmas[i])

		if err != nil {
			log.Println("Bad pragma: " + pragmas[i])
			return err
		}
	}
	return nil
}
