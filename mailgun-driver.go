// Copyright 2023 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package mailing

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailGunConfig struct {
	Domain string // your-domain.com
	APIKey string // your api key
}

type MailGunDriver struct {
	config         *MailGunConfig
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

var initiateMainGunSend = func(from string, rcpts []string, message []byte, d Driver) error {
	mgDriver := d.(*MailGunDriver)
	mg := mailgun.NewMailgun(mgDriver.config.Domain, mgDriver.config.APIKey)
	var m *mailgun.Message
	if mgDriver.htmlBody != "" {
		m = mg.NewMessage(
			from,
			mgDriver.subject,
			"",
			rcpts...,
		)
		m.SetHtml(mgDriver.htmlBody)
	} else {
		m = mg.NewMessage(
			from,
			mgDriver.subject,
			mgDriver.plainTextBody,
			rcpts...,
		)
	}
	if len(mgDriver.attachments) != 0 {
		for _, v := range mgDriver.attachments {
			m.AddAttachment(v.Path)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, _, err := mg.Send(ctx, m)
	if err != nil {
		return errors.New(fmt.Sprintf(" error calling mg.Send(): %v", err.Error()))
	}
	return nil
}

func initiateMainGun(config *SendGridConfig) *SendGridDriver {
	s := &SendGridDriver{
		config:         config,
		messageBuilder: newMessageBuilder(),
		htmlBody:       "",
		plainTextBody:  "",
		initiateSend:   initiateMainGunSend,
	}

	return s
}

func (m *MailGunDriver) SetFrom(from mail.Address) error {
	m.from = from
	return nil
}

func (m *MailGunDriver) SetTo(toList []mail.Address) error {
	m.toList = toList
	return nil
}

func (m *MailGunDriver) SetCC(ccList []mail.Address) error {
	m.ccList = ccList
	return nil
}
func (m *MailGunDriver) SetBCC(bccList []mail.Address) error {
	m.bccList = bccList
	return nil
}
func (m *MailGunDriver) SetSubject(Subject string) error {
	m.subject = Subject
	return nil
}
func (m *MailGunDriver) SetHTMLBody(body string) error {
	m.htmlBody = body
	return nil
}
func (m *MailGunDriver) SetPlainTextBody(body string) error {
	m.plainTextBody = body
	return nil
}
func (m *MailGunDriver) SetAttachments(attachments []Attachment) error {
	m.attachments = attachments
	return nil
}

func (m *MailGunDriver) Send() error {
	// prepare the message
	m.messageBuilder.setSubject(m.subject)
	if m.htmlBody != "" {
		m.messageBuilder.setHTMLBody(m.htmlBody)
	} else {
		m.messageBuilder.setPlainTextBody(m.plainTextBody)
	}
	m.messageBuilder.setFrom(m.from)
	m.messageBuilder.setToList(m.toList)
	m.messageBuilder.setCCList(m.ccList)
	m.messageBuilder.setAttachments(m.attachments)
	message := m.messageBuilder.build()

	// "to" and "cc" message sending
	var rcpts []string
	for _, v := range m.toList {
		rcpts = append(rcpts, v.String())
	}
	for _, v := range m.ccList {
		rcpts = append(rcpts, v.String())
	}
	from := m.from.String()
	err := m.initiateSend(from, rcpts, message, m)
	if err != nil {
		return errors.New(fmt.Sprintf("error calling m.initiateSend(): %v", err.Error()))
	}

	// send to bcc
	for _, v := range m.bccList {
		err = m.initiateSend(from, []string{v.String()}, message, m)
		if err != nil {
			return errors.New(fmt.Sprintf("error calling m.initiateSend(): %v", err.Error()))
		}
	}
	m.resetDriverProps()
	return nil
}

func (m *MailGunDriver) resetDriverProps() {
	m.subject = ""
	m.htmlBody = ""
	m.plainTextBody = ""
}
