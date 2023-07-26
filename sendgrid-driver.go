// Copyright 2023 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package mailing

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"os"

	"github.com/sendgrid/sendgrid-go"
	sgmail "github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridConfig struct {
	Host     string // "https://api.sendgrid.com"
	Endpoint string // "/v3/mail/send"
	ApiKey   string // SENDGRID_API_KEY
}

type SendGridDriver struct {
	config         *SendGridConfig
	messageBuilder *messageBuilder
	from           mail.Address
	toList         []mail.Address
	ccList         []mail.Address
	bccList        []mail.Address
	subject        string
	htmlBody       string
	plainTextBody  string
	attachments    []Attachment
	initiateSend   func(from string, rcpts []string, message []byte, conf Driver) error
}

var initiateSendGridSend = func(from string, rcpts []string, message []byte, d Driver) error {
	sgDriver := d.(*SendGridDriver)
	m := sgmail.NewV3Mail()
	fromEmail := sgmail.NewEmail(sgDriver.from.Name, sgDriver.from.Address)
	m.SetFrom(fromEmail)
	m.Subject = sgDriver.subject

	p := sgmail.NewPersonalization()
	if len(sgDriver.toList) != 0 {
		var tos []*sgmail.Email
		for _, v := range sgDriver.toList {
			tos = append(tos, sgmail.NewEmail(v.Name, v.Address))
		}
		p.AddTos(tos...)
	}
	if len(sgDriver.ccList) != 0 {
		var ccs []*sgmail.Email
		for _, v := range sgDriver.ccList {
			ccs = append(ccs, sgmail.NewEmail(v.Name, v.Address))
		}
		p.AddCCs(ccs...)
	}
	if len(sgDriver.bccList) != 0 {
		var bccs []*sgmail.Email
		for _, v := range sgDriver.bccList {
			bccs = append(bccs, sgmail.NewEmail(v.Name, v.Address))
		}
		p.AddBCCs(bccs...)
	}
	m.AddPersonalizations(p)
	if sgDriver.plainTextBody != "" {
		c := sgmail.NewContent("text/plain", sgDriver.plainTextBody)
		m.AddContent(c)
	}
	if sgDriver.htmlBody != "" {
		c := sgmail.NewContent("text/html", sgDriver.htmlBody)
		m.AddContent(c)
	}

	var a *sgmail.Attachment
	var attachementContent []byte
	var err error
	for _, v := range sgDriver.attachments {
		attachementContent, err = os.ReadFile(v.Path)
		if err != nil {
			return err
		}

		encodedAttachmentbuf := base64.StdEncoding.EncodeToString([]byte(attachementContent))
		a = sgmail.NewAttachment()
		a.SetContent(encodedAttachmentbuf)

		a.SetType(http.DetectContentType(attachementContent))

		a.SetFilename(v.Name)

		a.SetDisposition("attachment")

		m.AddAttachment(a)
	}

	requestBody := sgmail.GetRequestBody(m)

	request := sendgrid.GetRequest(sgDriver.config.ApiKey, sgDriver.config.Endpoint, sgDriver.config.Host)
	request.Method = "POST"
	var Body = requestBody
	request.Body = Body
	_, err = sendgrid.API(request)
	if err != nil {
		return err
	}

	return nil
}

func initiateSendGrid(config *SendGridConfig) *SendGridDriver {
	s := &SendGridDriver{
		config:         config,
		messageBuilder: newMessageBuilder(),
		htmlBody:       "",
		plainTextBody:  "",
		initiateSend:   initiateSendGridSend,
	}

	return s
}

func (s *SendGridDriver) SetFrom(from mail.Address) error {
	s.from = from
	return nil
}

func (s *SendGridDriver) SetTo(toList []mail.Address) error {
	s.toList = toList
	return nil
}

func (s *SendGridDriver) SetCC(ccList []mail.Address) error {
	s.ccList = ccList
	return nil
}
func (s *SendGridDriver) SetBCC(bccList []mail.Address) error {
	s.bccList = bccList
	return nil
}
func (s *SendGridDriver) SetSubject(Subject string) error {
	s.subject = Subject
	return nil
}
func (s *SendGridDriver) SetHTMLBody(body string) error {
	s.htmlBody = body
	return nil
}
func (s *SendGridDriver) SetPlainTextBody(body string) error {
	s.plainTextBody = body
	return nil
}
func (s *SendGridDriver) SetAttachments(attachments []Attachment) error {
	s.attachments = attachments
	return nil
}

func (s *SendGridDriver) Send() error {
	// prepare the message
	s.messageBuilder.setSubject(s.subject)
	if s.htmlBody != "" {
		s.messageBuilder.setHTMLBody(s.htmlBody)
	} else {
		s.messageBuilder.setPlainTextBody(s.plainTextBody)
	}
	s.messageBuilder.setFrom(s.from)
	s.messageBuilder.setToList(s.toList)
	s.messageBuilder.setCCList(s.ccList)
	s.messageBuilder.setAttachments(s.attachments)
	message := s.messageBuilder.build()

	// "to" and "cc" message sending
	var rcpts []string
	for _, v := range s.toList {
		rcpts = append(rcpts, v.String())
	}
	for _, v := range s.ccList {
		rcpts = append(rcpts, v.String())
	}
	from := s.from.String()
	err := s.initiateSend(from, rcpts, message, s)
	if err != nil {
		return errors.New(fmt.Sprintf("error calling s.initiateSend(): %v", err.Error()))
	}

	// send to bcc
	for _, v := range s.bccList {
		err = s.initiateSend(from, []string{v.String()}, message, s)
		if err != nil {
			return errors.New(fmt.Sprintf("error calling s.initiateSend(): %v", err.Error()))
		}
	}
	s.resetDriverProps()
	return nil
}

func (s *SendGridDriver) resetDriverProps() {
	s.subject = ""
	s.htmlBody = ""
	s.plainTextBody = ""
}
