package main

import (
	"encoding/xml"
	lib "github.com/gnewton/pstxml2sqlite/pstxml2sqlitestructs"
	"log"
)

var idCounter int64 = 0
var dups map[string]bool

var chunk []*lib.Message = nil
var n int

func handleFeed(saver *MessageSaver, se xml.StartElement, decoder *xml.Decoder, messagesChannel chan []*lib.Message) {
	log.Println("idCounter", idCounter)
	log.Println("local=" + se.Name.Local)
	if chunk == nil {
		chunk = make([]*lib.Message, chunkSize)
		n = 0
	}
	/*if se.Name.Local == "filesource" && se.Name.Space == "" {
		var fileSource lib.Filesource
		decoder.DecodeElement(&fileSource, &se)
		writeFileSource(fileSource)
	} else */
	//log.Println("local=" + se.Name.Local)
	if se.Name.Local == "message" && se.Name.Space == "" {
		log.Println("****")
		idCounter++
		counter++
		countAll++
		var message lib.Message
		err := decoder.DecodeElement(&message, &se)
		if err != nil {
			log.Fatal(err)
		}
		message.Id = idCounter
		chunk[n] = &message
		n++
		log.Println("n=", n)
		if n == chunkSize {
			log.Println("@@@@@ Sending to channel")
			messagesChannel <- chunk
			chunk = nil
			n = 0
		}
	}
}

func writeFileSource(fileSource lib.Filesource) {
	log.Println("filesource:", fileSource)
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
