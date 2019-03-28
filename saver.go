package main

import (
	//"crypto/sha256"
	"database/sql"
	//"fmt"
	lib "github.com/gnewton/pstxml2sqlite/pstxml2sqlitestructs"
	"log"
	"strconv"
	"sync"
)

type MessageSaver struct {
	messageStmt, recipStmt, attachStmt *sql.Stmt
	db2                                *sql.DB
	tx2                                *sql.Tx
	messagesChannel                    chan []*lib.Message
	transactionMux, dupMux             *sync.Mutex
}

var txCounter = 0
var txSizeCounter int64 = 0

const messageSql = "INSERT INTO messages(attr_folder_depth,attr_folders_path,attr_importance,attr_in_reply_to_id,attr_internet_article_number,attr_num_attachments,attr_priority,attr_response_requested,attr_sensitivity,body_raw,\"from\",from_name,has_attachments,id,is_cc_me,is_forwarded,is_from_me,is_message_recip_me,is_message_to_me,is_read,is_replied,is_response_requested,is_resent,is_submitted,is_unsent,is_unmodified,message_id,num_recipients,received ,return_path,subject,sha256) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

const attachmentSql = "INSERT INTO attachments(id,message_id,attr_attach_type,attr_attachment_content_disposition,attr_filename,attr_mime,size,raw_content,raw_size,extracted_text,sha256_hex) VALUES (?,?,?,?,?,?,?,?,?,?,?)"

const recipSql = "INSERT INTO recipients(id,message_id,attr_email,attr_mapi,attr_name,attr_smtp) values (?,?,?,?,?,?)"

const fsSql = "INSERT INTO filesources(fid,filename) values (?,?)"

func newStatement(tx *sql.Tx, sql string) (*sql.Stmt, error) {
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Println("Bad sql:", sql)
	}
	return stmt, err

}

//func (saver *MessageSaver) save(message *lib.Message) {
func (saver *MessageSaver) run(wg *sync.WaitGroup) {
	defer wg.Done()

	for messages := range saver.messagesChannel {
		log.Println("***** Size chunk:", len(messages))
		for i, _ := range messages {
			message := messages[i]

			log.Println("")
			log.Println("---------------------------------------------")
			log.Println("countAlll2", countAll2)
			log.Println("count2", countAll)
			log.Println("txCounter", txCounter)
			if message == nil {
				continue
			}
			//log.Println(i, message.Id)

			countAll2++
			fixMessageFields(message)
			if message.Recipients != nil {
				prepareRecipients(message.Recipients.Recipient, idCounter)
			}

			if message.Attachments != nil {
				prepareAttachments(message.Attachments.Attachment, idCounter)
			}
			saver.transactionMux.Lock()

			txSizeCounter += message.Length()
			if txCounter == TxSize || txSizeCounter > TxLengthSize {
				if txSizeCounter > TxLengthSize {
					log.Println("$$$$$$$$$$$$$$$$$$$ txSizeCounter > TxLengthSize", txCounter, txSizeCounter-message.Length(), TxLengthSize)
				}
				saver.tx2.Commit()
				log.Println("COMMIT")
				txCounter = 0
				txSizeCounter = 0
				var err error
				saver.tx2, err = saver.db2.Begin()
				if err != nil {
					log.Fatal(err)
				}
				saver.messageStmt, err = newStatement(saver.tx2, messageSql)
				if err != nil {
					log.Fatal(err)
				}
				saver.attachStmt, err = newStatement(saver.tx2, attachmentSql)
				if err != nil {
					log.Fatal(err)
				}
				saver.recipStmt, err = newStatement(saver.tx2, recipSql)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("countAll=", countAll)
				log.Println("countAll2=", countAll2)

			}
			newSave(message, saver.messageStmt, saver.attachStmt, saver.recipStmt)
			txCounter++
			log.Println("messageSize=", message.Length())
			txSizeCounter += message.Length()
			saver.transactionMux.Unlock()
		}

	}
}
func newSave(message *lib.Message, messageStmt, attachStmt, recipStmt *sql.Stmt) {

	//_, err := stmt.Exec(message.Body)
	_, err := messageStmt.Exec(message.AttrFolderDepth, message.AttrFoldersPath, message.AttrImportance, message.AttrInReplyToId, message.AttrInternetArticleNumber, message.AttrNumAttachments, message.AttrPriority, message.AttrResponseRequested, message.AttrSensitivity, message.BodyRaw, message.From, message.From_name, message.HasAttachments, message.Id, message.IsCcMe, message.IsForwarded, message.IsFromMe, message.IsMessageRecipMe, message.IsMessageToMe, message.IsRead, message.IsReplied, message.IsResponseRequested, message.IsResent, message.IsSubmitted, message.IsUnsent, message.IsUnmodified, message.Message_id, message.Num_recipients, message.Received, message.Return_path, message.Subject, message.SHA256)

	if err != nil {
		log.Fatal(err)
	}

	if message.Attachments != nil {
		atts := message.Attachments.Attachment
		for i := 0; i < len(atts); i++ {
			att := atts[i]
			att.Id = attachCounter
			att.MessageId = message.Id

			var extractedText string
			if att.ContentTextExtracted == nil || att.ContentTextExtracted.Content == nil || att.ContentTextExtracted.Content.Data == "" {
				extractedText = ""
			} else {
				extractedText = att.ContentTextExtracted.Content.Data
			}

			_, err := attachStmt.Exec(att.Id, att.MessageId, att.AttrAttachType, att.AttrAttachmentContentDisposition, att.AttrFilename, att.AttrMime, att.Size, att.RawContent, att.RawSize, NewNullString(extractedText), att.Sha256Hex)
			if err != nil {
				log.Fatal(err)
			}
			attachCounter++
		}
	}

	if message.Recipients != nil {
		recips := message.Recipients.Recipient
		for i := 0; i < len(recips); i++ {
			recip := recips[i]
			recip.Id = recipCounter
			recip.MessageId = message.Id
			_, err := recipStmt.Exec(recip.Id, recip.MessageId, recip.AttrEmail, recip.AttrMapi, recip.AttrName, recip.AttrSmtp)
			if err != nil {
				log.Fatal(err)
			}
			recipCounter++
		}
	}

}

// From: https://stackoverflow.com/a/40268372
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
func prepareRecipients(recips []*lib.Recipient, messageId int64) {
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
func prepareAttachments(attachments []*lib.Attachment, messageId int64) {
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
		if len(attach.Content.Data) > 0 {
			var err error
			attach.RawContent, err = decodeBase64(attach.Content.Data)
			if err != nil {
				log.Fatal(err)
			}
			attach.RawSize = len(attach.RawContent)
		}
	}
}

func fixMessageFields(mes *lib.Message) {
	t, _ := string2date(mes.OrigReceived)
	mes.Received = *t

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
	tmpBodyRaw, err := decodeBase64(mes.Body.Data)
	if err != nil {
		log.Println(mes.Body)
		log.Fatal(err)
	}
	mes.BodyRaw = string(tmpBodyRaw)

}
