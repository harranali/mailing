package mailing

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/mail"
	"os"
	"strings"
)

type messageBuilder struct {
	subject       string
	htmlBody      string
	plainTextBody string
	from          string
	toList        []string
	ccList        []string
	attachments   []Attachment
}

func newMessageBuilder() *messageBuilder {
	return &messageBuilder{
		subject:       "",
		htmlBody:      "",
		plainTextBody: "",
	}
}

func (m *messageBuilder) setSubject(subject string) *messageBuilder {
	m.subject = subject
	return m
}

func (m *messageBuilder) setHTMLBody(body string) *messageBuilder {
	m.htmlBody = body
	return m
}

func (m *messageBuilder) setPlainTextBody(body string) *messageBuilder {
	m.plainTextBody = body
	return m
}

func (m *messageBuilder) setFrom(from mail.Address) *messageBuilder {
	m.from = from.String()
	return m
}

func (m *messageBuilder) setToList(toList []mail.Address) *messageBuilder {
	var toListStr []string
	for _, v := range toList {
		toListStr = append(toListStr, v.String())
	}
	m.toList = toListStr
	return m
}

func (m *messageBuilder) setCCList(ccList []mail.Address) *messageBuilder {
	var ccListStr []string
	for _, v := range ccList {
		ccListStr = append(ccListStr, v.String())
	}
	m.ccList = ccListStr
	return m
}

func (m *messageBuilder) setAttachments(attachments []Attachment) *messageBuilder {
	m.attachments = attachments
	return m
}

func (m *messageBuilder) build() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("From: %s\r\n", m.from))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(m.toList, ";")))
	buf.WriteString(fmt.Sprintf("CC: %s\r\n", strings.Join(m.ccList, ";")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", m.subject))
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary))
	if m.htmlBody != "" {
		buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
		buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
		buf.WriteString(m.htmlBody)
	} else {
		buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
		buf.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
		buf.WriteString(m.plainTextBody)
	}
	if len(m.attachments) > 0 {
		for _, attachment := range m.attachments {
			file, err := os.Open(attachment.Path)
			fileContent, err := io.ReadAll(file)
			if err != nil {
				panic(err.Error())
			}
			if err != nil {
				panic(err.Error())
			}
			buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: \"%s\"\r\n", http.DetectContentType(fileContent)))
			buf.WriteString("Content-Transfer-Encoding: base64\r\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", attachment.Name))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(fileContent)))
			base64.StdEncoding.Encode(b, fileContent)
			buf.Write(b)
			file.Close()
		}
	}
	buf.WriteString(fmt.Sprintf("\r\n--%s--\r\n", boundary))
	return buf.Bytes()
}
