// Copyright 2023 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package mailing

import (
	"crypto/tls"
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
	conn           *tls.Conn
	client         *smtp.Client
	messageBuilder *messageBuilder
	from           mail.Address
	toList         []mail.Address
	ccList         []mail.Address
	bccList        []mail.Address
	subject        string
	htmlBody       string
	plainTextBody  string
	attachments    []Attachment
}

func initiateSMTP(config *SMTPConfig) *smtpDriver {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port), &config.TLSConfig)
	if err != nil {
		panic(err.Error())
	}
	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		panic(err.Error())
	}
	err = client.Auth(smtp.PlainAuth("", config.Username, config.Password, config.Host))
	if err != nil {
		panic(err.Error())
	}
	return &smtpDriver{
		conn:           conn,
		client:         client,
		messageBuilder: newMessageBuilder(),
		htmlBody:       "",
		plainTextBody:  "",
	}
}

func (s *smtpDriver) SetFrom(from mail.Address) *smtpDriver {
	s.from = from
	return s
}

func (s *smtpDriver) SetTo(toList []mail.Address) *smtpDriver {
	s.toList = toList
	return s
}

func (s *smtpDriver) SetCC(ccList []mail.Address) *smtpDriver {
	s.ccList = ccList
	return s
}
func (s *smtpDriver) SetBCC(bccList []mail.Address) *smtpDriver {
	s.bccList = bccList
	return s
}
func (s *smtpDriver) SetSubject(Subject string) *smtpDriver {
	s.subject = Subject
	return s
}
func (s *smtpDriver) SetHTMLBody(body string) *smtpDriver {
	s.htmlBody = body
	return s
}
func (s *smtpDriver) SetPlainTextBody(body string) *smtpDriver {
	s.plainTextBody = body
	return s
}
func (s *smtpDriver) SetAttachments(attachments []Attachment) *smtpDriver {
	s.attachments = attachments
	return s
}

func (s *smtpDriver) Send() error {
	s.client.Mail(prepareAddressString(s.from))
	for _, emailAddress := range s.toList {
		s.client.Rcpt(prepareAddressString(emailAddress))
	}
	for _, emailAddress := range s.ccList {
		s.client.Rcpt(prepareAddressString(emailAddress))
	}
	for _, emailAddress := range s.bccList {
		s.client.Rcpt(prepareAddressString(emailAddress))
	}
	writer, err := s.client.Data()
	if err != nil {
		return err
	}

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
	_, err = writer.Write(message)
	if err != nil {
		return err
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	err = s.client.Quit()
	if err != nil {
		return err
	}
	err = s.client.Close()
	if err != nil {
		return err
	}
	err = s.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
