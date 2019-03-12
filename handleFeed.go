package main

import (
	//"crypto/sha256"
	"database/sql"
	"encoding/xml"
	//"fmt"
	lib "github.com/gnewton/pstxml2sqlite/pstxml2sqlitestructs"
	"github.com/jinzhu/gorm"
	"log"
	"strconv"
	"sync"
)

var idCounter int64 = 0
var dups map[string]bool

var chunk []*lib.Message = nil
var n int
var chunkSize = 50

func handleFeed(saver *MessageSaver, se xml.StartElement, decoder *xml.Decoder, c chan []*lib.Message) {
	if chunk == nil {
		chunk = make([]*lib.Message, chunkSize)
		n = 0
	}
	if se.Name.Local == "message" && se.Name.Space == "" {
		idCounter++
		//log.Println(countAll, counter)
		counter++
		countAll++
		var message lib.Message
		decoder.DecodeElement(&message, &se)
		message.Id = idCounter
		//saver.save(&message)
		chunk[n] = &message
		n++
		if n == chunkSize {
			//c <- &message
			c <- chunk
			chunk = nil
		}
	}
}

type MessageSaver struct {
	stmt    *sql.Stmt
	db, tx  *gorm.DB
	db2     *sql.DB
	tx2     *sql.Tx
	c       chan []*lib.Message
	counter int32
	mux     sync.Mutex
}

const messageSql = "insert into messages(attr_folder_depth,attr_folders_path,attr_importance,attr_in_reply_to_id,attr_internet_article_number,attr_num_attachments,attr_priority,attr_response_requested,attr_sensitivity,body_raw,\"from\",from_name,has_attachments,id,is_cc_me,is_forwarded,is_from_me,is_message_recip_me,is_message_to_me,is_read,is_replied,is_response_requested,is_resent,is_submitted,is_unsent,is_unmodified,message_id,num_recipients,received ,return_path,subject,sha256) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

func newStatement(tx *sql.Tx) (*sql.Stmt, error) {
	//return tx.Prepare("insert into messages(body_raw) values(?)")
	return tx.Prepare(messageSql)

}

//func (saver *MessageSaver) save(message *lib.Message) {
func (saver *MessageSaver) run(wg *sync.WaitGroup) {
	defer wg.Done()

	for messages := range saver.c {
		for i, _ := range messages {
			message := messages[i]
			//log.Println(i, message.Id)
			fixMessageFields(message)

			//tmp := []byte(message.Received.String() + message.From + message.AttrInternetArticleNumber + message.BodyRaw)

			//message.SHA256 = fmt.Sprintf("%x", sha256.Sum256(tmp))
			//saver.mux.Lock()
			//if _, ok := dups[message.SHA256]; ok {
			//	saver.mux.Unlock()
			//	return
			//}

			dups[message.SHA256] = true
			//saver.mux.Unlock()

			if message.Recipients != nil {
				prepareRecipients(message.Recipients.Recipient, saver.db, idCounter)
			}

			if message.Attachments != nil {
				prepareAttachments(message.Attachments.Attachment, saver.db, idCounter)
			}
			//saveAll(message, saver.db)
			saver.mux.Lock()
			newSave(message, saver.stmt)
			saver.counter++
			if saver.counter == 1000 {
				saver.tx2.Commit()
				saver.counter = 0
				var err error
				saver.tx2, err = saver.db2.Begin()
				if err != nil {
					log.Fatal(err)
				}
				saver.db = saver.tx
				saver.stmt, err = newStatement(saver.tx2)
				if err != nil {
					log.Fatal(err)
				}
				log.Println(countAll)
			}

			saver.mux.Unlock()
		}

	}
}
func newSave(message *lib.Message, stmt *sql.Stmt) {
	//log.Println(message.Id)
	//"attr_folder_depth" varchar(255),"attr_folders_path" varchar(255),"attr_importance" varchar(255),"attr_in_reply_to_id" varchar(255),"attr_internet_article_number" varchar(255),"attr_num_attachments" varchar(255),"attr_priority" varchar(255),"attr_response_requested" varchar(255),"attr_sensitivity" varchar(255),"body_raw" text,"from" varchar(255),"from_name" varchar(255),"has_attachments" bool,"id" integer primary key autoincrement,"is_cc_me" bool,"is_forwarded" bool,"is_from_me" bool,"is_message_recip_me" bool,"is_message_to_me" bool,"is_read" bool,"is_replied" bool,"is_response_requested" bool,"is_resent" bool,"is_submitted" bool,"is_unsent" bool,"is_unmodified" bool,"message_id" varchar(255),"num_recipients" varchar(255),"received" datetime,"return_path" varchar(255),"subject" varchar(255),"sha256" varchar(64) );

	//_, err := stmt.Exec(message.Body)
	_, err := stmt.Exec(message.AttrFolderDepth, message.AttrFoldersPath, message.AttrImportance, message.AttrInReplyToId, message.AttrInternetArticleNumber, message.AttrNumAttachments, message.AttrPriority, message.AttrResponseRequested, message.AttrSensitivity, message.BodyRaw, message.From, message.From_name, message.HasAttachments, message.Id, message.IsCcMe, message.IsForwarded, message.IsFromMe, message.IsMessageRecipMe, message.IsMessageToMe, message.IsRead, message.IsReplied, message.IsResponseRequested, message.IsResent, message.IsSubmitted, message.IsUnsent, message.IsUnmodified, message.Message_id, message.Num_recipients, message.Received, message.Return_path, message.Subject, message.SHA256)

	if err != nil {
		log.Fatal(err)
	}
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
