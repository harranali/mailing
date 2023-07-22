// Copyright 2023 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package mailing

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/SparkPost/gosparkpost"
)

type SparkPostConfig struct {
	BaseUrl    string // example: https://api.sparkpost.com
	ApiKey     string // SPARKPOST_API_KEY
	ApiVersion int    // example: 1
}

type SparkPostDriver struct {
	config         *SparkPostConfig
	messageBuilder *messageBuilder
	from           mail.Address
	toList         []mail.Address
	ccList         []mail.Address
	bccList        []mail.Address
	subject        string
	htmlBody       string
	plainTextBody  string
	attachments    []Attachment
	initiateSend   func(from string, rcpts []string, message []byte, conf Driver) (id string, err error)
}

var initiateSend = func(from string, rcpts []string, message []byte, d Driver) (id string, err error) {
	spDriv := d.(*SparkPostDriver)
	conf := spDriv.config
	cfg := &gosparkpost.Config{
		BaseUrl:    conf.BaseUrl,
		ApiKey:     conf.ApiKey,
		ApiVersion: conf.ApiVersion,
	}
	var client gosparkpost.Client
	err = client.Init(cfg)
	if err != nil {
		return "", errors.New(fmt.Sprintf("SparkPost client init failed: %s\n", err))
	}

	// create the content
	content := gosparkpost.Content{
		From:    from,
		Subject: spDriv.subject,
	}
	// the body
	if spDriv.htmlBody != "" {
		content.HTML = spDriv.htmlBody
	} else {
		content.Text = spDriv.plainTextBody
	}
	// headers
	var headers map[string]string
	// from
	headers["From"] = spDriv.from.String()
	// the cc
	if len(spDriv.messageBuilder.ccList) > 0 {
		headers["Cc"] = strings.Join(spDriv.messageBuilder.ccList, ",")
	}
	// add the headers
	if len(headers) > 0 {
		content.Headers = headers
	}
	// Create transmission
	tx := &gosparkpost.Transmission{
		Recipients: rcpts,
		Content:    content,
	}
	id, _, err = client.Send(tx)
	if err != nil {
		return "", err
	}

	return id, nil
}

func initiateSparkPost(config *SparkPostConfig) *SparkPostDriver {
	s := &SparkPostDriver{
		config:         config,
		messageBuilder: newMessageBuilder(),
		htmlBody:       "",
		plainTextBody:  "",
		initiateSend:   initiateSend,
	}

	return s
}

func (s *SparkPostDriver) SetFrom(from mail.Address) error {
	s.from = from
	return nil
}

func (s *SparkPostDriver) SetTo(toList []mail.Address) error {
	s.toList = toList
	return nil
}

func (s *SparkPostDriver) SetCC(ccList []mail.Address) error {
	s.ccList = ccList
	return nil
}
func (s *SparkPostDriver) SetBCC(bccList []mail.Address) error {
	s.bccList = bccList
	return nil
}
func (s *SparkPostDriver) SetSubject(Subject string) error {
	s.subject = Subject
	return nil
}
func (s *SparkPostDriver) SetHTMLBody(body string) error {
	s.htmlBody = body
	return nil
}
func (s *SparkPostDriver) SetPlainTextBody(body string) error {
	s.plainTextBody = body
	return nil
}
func (s *SparkPostDriver) SetAttachments(attachments []Attachment) error {
	s.attachments = attachments
	return nil
}

func (s *SparkPostDriver) Send() error {
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
	_, err := s.initiateSend(from, rcpts, message, s)
	if err != nil {
		return errors.New(fmt.Sprintf("error calling s.initiateSend(): %v", err.Error()))
	}

	// send to bcc
	for _, v := range s.bccList {
		_, err = s.initiateSend(from, []string{v.String()}, message, s)
		if err != nil {
			return errors.New(fmt.Sprintf("error calling s.initiateSend(): %v", err.Error()))
		}
	}
	s.resetDriverProps()
	return nil
}

func (s *SparkPostDriver) resetDriverProps() {
	s.subject = ""
	s.htmlBody = ""
	s.plainTextBody = ""
}
