package main

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	lib "github.com/gnewton/pstxml2sqlite/pstxml2sqlitestructs"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
)

var idCounter int64 = 0
var dups map[string]bool

func handleFeed(saver *MessageSaver, se xml.StartElement, decoder *xml.Decoder, db *gorm.DB) {
	if se.Name.Local == "message" && se.Name.Space == "" {
		idCounter++
		//log.Println(countAll, counter)
		counter++
		countAll++
		var message lib.Message
		decoder.DecodeElement(&message, &se)
		message.Id = idCounter
		saver.save(&message)

	}
}

type MessageSaver struct {
	db *gorm.DB
	//c  chan (*lib.Message)
}

func (saver *MessageSaver) save(message *lib.Message) {
	///// push to channel HERE-----
	fixMessageFields(message)

	tmp := []byte(message.Received.String() + message.From + message.AttrInternetArticleNumber + message.BodyRaw)

	message.SHA256 = fmt.Sprintf("%x", sha256.Sum256(tmp))

	if _, ok := dups[message.SHA256]; ok {
		return
	}
	//log.Println(message.SHA256)
	dups[message.SHA256] = true

	if message.Recipients != nil {
		prepareRecipients(message.Recipients.Recipient, saver.db, idCounter)
	}

	if message.Attachments != nil {
		prepareAttachments(message.Attachments.Attachment, saver.db, idCounter)
	}
	saveAll(message, saver.db)

}

func saveAll(message *lib.Message, db *gorm.DB) {
	db.Create(&message)

	//Recipients
	if message.Recipients != nil {
		for i := 0; i < len(message.Recipients.Recipient); i++ {
			//recip := message.Recipients.Recipient[i]
			_ = message.Recipients.Recipient[i]
			//db.Create(&recip)
		}
	}

	//Attachments
	if message.Attachments != nil {
		for i := 0; i < len(message.Attachments.Attachment); i++ {
			//attach := message.Attachments.Attachment[i]
			_ = message.Attachments.Attachment[i]
			//db.Create(&attach)
		}
	}
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

func prepareRecipients(recips []*lib.Recipient, db *gorm.DB, messageId int64) {
	for i := 0; i < len(recips); i++ {
		recip := recips[i]
		recip.MessageId = messageId
	}
}

// We realy only want attachmentType=1, which has the actual file. The others are references.
// From: https://docs.microsoft.com/en-us/office/vba/api/outlook.olattachmenttype
// Name 	  Value 	Description
// olByReference  4 	This value is no longer supported since Microsoft Outlook 2007. Use olByValue to attach a copy of a file in the file system.
// olByValue 	  1 	The attachment is a copy of the original file and can be accessed even if the original file is removed.
// olEmbeddeditem 5 	The attachment is an Outlook message format file (.msg) and is a copy of the original message.
// olOLE 	  6 	The attachment is an OLE document.
//
func prepareAttachments(attachments []*lib.Attachment, db *gorm.DB, messageId int64) {
	for i := 0; i < len(attachments); i++ {
		attach := attachments[i]
		//log.Println("Filename: " + attach.AttrFilename)

		attach.MessageId = messageId
		attach.Size = 0
		if len(attach.AttrSize) > 0 {
			size, err := strconv.ParseInt(attach.AttrSize, 10, 64)
			if err != nil {
				log.Fatal(err)
			} else {
				attach.Size = size
			}
		}
		if len(attach.Content.Text) > 0 {
			var err error
			attach.RawContent, err = decodeBase64(attach.Content.Text)
			if err != nil {
				log.Fatal(err)
			}
			attach.RawSize = len(attach.RawContent)
		}
	}
}

func fixMessageFields(mes *lib.Message) {
	mes.Received, _ = string2date(mes.OrigReceived)

	// booleans
	mes.IsCcMe = string2bool(mes.AttrCcMe)
	mes.IsForwarded = string2bool(mes.AttrForwarded)
	//mes.IsFromMe = string2bool(mes.AttrFromMe)
	//mes.IsMessageRecipMe = string2bool(mes.AttrMessageRecipMe)
	//mes.IsMessageToMe = string2bool(mes.AttrMessageToMe)
	mes.IsRead = string2bool(mes.AttrRead)
	mes.IsReplied = string2bool(mes.AttrReplied)
	mes.IsResponseRequested = string2bool(mes.AttrResponseRequested)
	mes.IsResent = string2bool(mes.AttrResent)
	mes.IsSubmitted = string2bool(mes.AttrSubmitted)
	mes.IsUnmodified = string2bool(mes.AttrUnmodified)
	mes.IsUnsent = string2bool(mes.AttrUnsent)

	//base64 encoded
	tmpBodyRaw, err := decodeBase64(mes.Body)
	if err != nil {
		log.Fatal(err)
	}
	mes.BodyRaw = string(tmpBodyRaw)

}
