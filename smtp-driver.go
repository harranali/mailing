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
}

func initiateSMTP(config *SMTPConfig) *smtpDriver {
	return &smtpDriver{
		config:         config,
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
	err := s.initiateSend(rcpts, message)
	if err != nil {
		return errors.New(fmt.Sprintf("error calling s.initiateSend(): %v", err.Error()))
	}

	// send to bcc
	for _, v := range s.bccList {
		err = s.initiateSend([]string{v.String()}, message)
		if err != nil {
			return errors.New(fmt.Sprintf("error calling s.initiateSend(): %v", err.Error()))
		}
	}
	return nil
}

func (s *smtpDriver) initiateSend(rcpts []string, message []byte) error {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port), &s.config.TLSConfig)
	if err != nil {
		return errors.New(fmt.Sprintf("error calling tls.Dial(): %v", err.Error()))
	}
	defer conn.Close()
	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return errors.New(fmt.Sprintf("error calling smtp.NewClient(): %v", err.Error()))
	}
	defer client.Close()
	err = client.Auth(smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host))
	if err != nil {
		return errors.New(fmt.Sprintf("error calling SMTP's client.Auth(): %v", err.Error()))
	}
	client.Mail(s.from.String())
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
