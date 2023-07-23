// Copyright 2023 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package mailing

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
)

type SMTPConfig struct {
	Host      string
	Port      int
	Username  string
	Password  string
	TLSConfig tls.Config
}

type smtpDriver struct {
	config         *SMTPConfig
	messageBuilder *messageBuilder
	from           mail.Address
	toList         []mail.Address
	ccList         []mail.Address
	bccList        []mail.Address
	subject        string
	htmlBody       string
	plainTextBody  string
	attachments    []Attachment
	initiateSend   func(from string, rcpts []string, message []byte, d Driver) error
}

var smtpInitiateSend = func(from string, rcpts []string, message []byte, d Driver) error {
	smtpDriv := d.(*smtpDriver)
	conf := smtpDriv.config
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", conf.Host, conf.Port), &conf.TLSConfig)
	if err != nil {
		return errors.New(fmt.Sprintf("error calling tls.Dial(): %v", err.Error()))
	}
	defer conn.Close()
	client, err := smtp.NewClient(conn, conf.Host)
	if err != nil {
		return errors.New(fmt.Sprintf("error calling smtp.NewClient(): %v", err.Error()))
	}
	defer client.Close()
	err = client.Auth(smtp.PlainAuth("", conf.Username, conf.Password, conf.Host))
	if err != nil {
		return errors.New(fmt.Sprintf("error calling SMTP's client.Auth(): %v", err.Error()))
	}
	client.Mail(from)
	for _, emailAddress := range rcpts {
		err = client.Rcpt(emailAddress)
		if err != nil {
			return errors.New(fmt.Sprintf("error calling rcpt(): %v", err.Error()))
		}
	}
	writer, err := client.Data()
	if err != nil {
		return errors.New(fmt.Sprintf("error calling data(): %v", err.Error()))
	}
	_, err = writer.Write(message)
	if err != nil {
		return errors.New(fmt.Sprintf("error calling writer.Close(): %v", err.Error()))
	}
	writer.Close()
	err = client.Quit()
	if err != nil {
		return errors.New(fmt.Sprintf("error quiting client: %v", err.Error()))
	}
	return nil
}

func initiateSMTP(config *SMTPConfig) *smtpDriver {
	s := &smtpDriver{
		config:         config,
		messageBuilder: newMessageBuilder(),
		htmlBody:       "",
		plainTextBody:  "",
		initiateSend:   smtpInitiateSend,
	}

	return s
}

func (s *smtpDriver) SetFrom(from mail.Address) error {
	s.from = from
	return nil
}

func (s *smtpDriver) SetTo(toList []mail.Address) error {
	s.toList = toList
	return nil
}

func (s *smtpDriver) SetCC(ccList []mail.Address) error {
	s.ccList = ccList
	return nil
}
func (s *smtpDriver) SetBCC(bccList []mail.Address) error {
	s.bccList = bccList
	return nil
}
func (s *smtpDriver) SetSubject(Subject string) error {
	s.subject = Subject
	return nil
}
func (s *smtpDriver) SetHTMLBody(body string) error {
	s.htmlBody = body
	return nil
}
func (s *smtpDriver) SetPlainTextBody(body string) error {
	s.plainTextBody = body
	return nil
}
func (s *smtpDriver) SetAttachments(attachments []Attachment) error {
	s.attachments = attachments
	return nil
}

func (s *smtpDriver) Send() error {
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

func (s *smtpDriver) resetDriverProps() {
	s.subject = ""
	s.htmlBody = ""
	s.plainTextBody = ""
}
