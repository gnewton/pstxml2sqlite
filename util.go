package main

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
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

func string2date(str string) (time.Time, error) {

	layout := "02-01-2006 06:04:05"
	t, err := time.Parse(layout, str)

	if err != nil {
		log.Println(err)
		log.Println(str)
		layout := "01-02-2006 06:04:05"
		t, err = time.Parse(layout, str)
		if err != nil {
			log.Println(err)
			log.Println(str)
		} else {
			log.Println("Parsed ok, witd DD-MM-YY")
			log.Println(str)
		}
	}
	return t, err
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
