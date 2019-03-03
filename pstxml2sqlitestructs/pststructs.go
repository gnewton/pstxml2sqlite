package pstxml2sqlitestructs

/////////////////////////////////////////////////////////////////
//Code generated by chidley https://github.com/gnewton/chidley //
/////////////////////////////////////////////////////////////////

import (
	"time"
)

type ChidleyRoot314159 struct {
	Messages *Messages `xml:"messages,omitempty" json:"messages,omitempty"`
}

type Attachment struct {
	MessageId                        int64
	AttrAttachType                   string   `xml:"attachType,attr"  json:",omitempty"`
	AttrAttachmentContentDisposition string   `xml:"attachmentContentDisposition,attr"  json:",omitempty"`
	AttrFilename                     string   `xml:"filename,attr"  json:",omitempty"`
	AttrMime                         string   `xml:"mime,attr"  json:",omitempty"`
	AttrSize                         string   `gorm:"-" xml:"size,attr"  json:",omitempty"`
	Size                             int64    `xml:"length,attr"  json:",omitempty"`
	Content                          *Content `xml:"content,omitempty" json:"content,omitempty"`
	//Base64Content                    string
	RawContent []byte
	RawSize    int
	Sha1Base64 string `xml:"sha1Base64,omitempty" json:"sha1Base64,omitempty"`
}

type Attachments struct {
	Attachment []*Attachment `xml:"attachment,omitempty" json:"attachment,omitempty"`
}

type Content struct {
	Text string `xml:",chardata" json:",omitempty"`
}

type Message struct {
	Attachments               *Attachments `xml:"attachments,omitempty" json:"attachments,omitempty"`
	AttrAttachmentsOrig       string       `gorm:"-" xml:"attachments,attr"  json:",omitempty"`
	AttrCcMe                  string       `gorm:"-" xml:"ccMe,attr"  json:",omitempty"`
	AttrFolderDepth           string       `xml:"folderDepth,attr"  json:",omitempty"`
	AttrFoldersPath           string       `xml:"foldersPath,attr"  json:",omitempty"`
	AttrForwarded             string       `gorm:"-" xml:"forwarded,attr"  json:",omitempty"`
	AttrImportance            string       `xml:"importance,attr"  json:",omitempty"`
	AttrInReplyToId           string       `xml:"inReplyToId,attr"  json:",omitempty"`
	AttrInternetArticleNumber string       `xml:"internetArticleNumber,attr"  json:",omitempty"`
	AttrNumAttachments        string       `xml:"numAttachments,attr"  json:",omitempty"`
	AttrPriority              string       `xml:"priority,attr"  json:",omitempty"`
	AttrResent                string       `gorm:"-" xml:"resent,attr"  json:",omitempty"`
	AttrRead                  string       `gorm:"-" xml:"read,attr"  json:",omitempty"`
	AttrReplied               string       `gorm:"-" xml:"replied,attr"  json:",omitempty"`
	AttrResponseRequested     string       `xml:"responseRequested,attr"  json:",omitempty"`
	AttrSensitivity           string       `xml:"sensitivity,attr"  json:",omitempty"`
	AttrSubmitted             string       `gorm:"-" xml:"submitted,attr"  json:",omitempty"`
	AttrUnmodified            string       `gorm:"-" xml:"unmodified,attr"  json:",omitempty"`
	AttrUnsent                string       `gorm:"-" xml:"unsent,attr"  json:",omitempty"`
	Body                      string       `gorm:"-" xml:"body,omitempty" json:"body,omitempty"`
	BodyRaw                   string       `gorm:"type:text" xml:"body_raw,omitempty" json:"body_raw,omitempty"`
	From                      string       `gorm:"index" xml:"from,omitempty" json:"from,omitempty"`
	From_name                 string       `gorm:"index" xml:"from_name,omitempty" json:"from_name,omitempty"`
	HasAttachments            bool         `gorm:"index" xml:"hattachments,attr"  json:",omitempty"`
	Id                        int64        `gorm:"index,PRIMARY_KEY"`
	IsCcMe                    bool         `gorm:"index" xml:"isCcMe,attr"  json:",omitempty"`
	IsForwarded               bool         `gorm:"index" xml:"isForwarded,attr"  json:",omitempty"`
	IsFromMe                  bool         `gorm:"index" xml:"fromMe,attr"  json:",omitempty"`
	IsMessageRecipMe          bool         `xml:"messageRecipMe,attr"  json:",omitempty"`
	IsMessageToMe             bool         `xml:"messageToMe,attr"  json:",omitempty"`
	IsRead                    bool         `gorm:"index" xml:"isRead,attr"  json:",omitempty"`
	IsReplied                 bool         `xml:"isReplied,attr"  json:",omitempty"`
	IsResponseRequested       bool         `xml:"isReplyRequested,attr"  json:",omitempty"`
	IsResent                  bool         `xml:"isSent,attr"  json:",omitempty"`
	IsSubmitted               bool         `xml:"isSubmitted,attr"  json:",omitempty"`
	IsUnsent                  bool         `xml:"isUnSent,attr"  json:",omitempty"`
	IsUnmodified              bool         `xml:"isUnModified,attr"  json:",omitempty"`
	Message_id                string       `gorm:"index" xml:"message_id,omitempty" json:"message_id,omitempty"`
	Num_recipients            string       `xml:"num_recipients,omitempty" json:"num_recipients,omitempty"`
	OrigReceived              string       `gorm:"-" xml:"received,omitempty" json:"received,omitempty"`
	Received                  time.Time    `gorm:"index" xml:"dateReceived,omitempty" json:"dateReceived,omitempty"`
	Recipients                *Recipients  `xml:"recipients,omitempty" json:"recipients,omitempty"`
	Return_path               string       `xml:"return_path,omitempty" json:"return_path,omitempty"`
	Subject                   string       `xml:"subject,omitempty" json:"subject,omitempty"`
}

type Body struct {
	// base64 encoded
	Text string `xml:",chardata" json:",omitempty"`
}

type Messages struct {
	Message []*Message `xml:"message,omitempty" json:"message,omitempty"`
	Meta    *Meta      `xml:"meta,omitempty" json:"meta,omitempty"`
}

type Meta struct {
	AttrKey   string `xml:"key,attr"  json:",omitempty"`
	AttrValue string `xml:"value,attr"  json:",omitempty"`
}

type Recipient struct {
	MessageId int64  `gorm:"index"`
	AttrEmail string `gorm:"index" xml:"email,attr"  json:",omitempty"`
	AttrMapi  string `gorm:"index" xml:"mapi,attr"  json:",omitempty"`
	AttrName  string `gorm:"index" xml:"name,attr"  json:",omitempty"`
	AttrSmtp  string `gorm:"index" xml:"smtp,attr"  json:",omitempty"`
}

type Recipients struct {
	Recipient []*Recipient `xml:"recipient,omitempty" json:"recipient,omitempty"`
}

///////////////////////////
