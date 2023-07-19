// Copyright 2023 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package mailing

import (
	"crypto/tls"
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
	client *smtp.Client
}

func initiateSMTP(config *SMTPConfig) *smtpDriver {
	return nil
}

func (s *smtpDriver) SetFrom(from mail.Address) {
	// TODO implement
}

func (s *smtpDriver) SetTo(toList []mail.Address) {
	// TODO implement
}

func (s *smtpDriver) SetCC(ccList []mail.Address) {
	// TODO implement
}
func (s *smtpDriver) SetBCC(bccList []mail.Address) {
	// TODO implement
}
func (s *smtpDriver) SetSubject(Subject string) {
	// TODO implement
}
func (s *smtpDriver) SetHTMLBody(body string) {
	// TODO implement
}
func (s *smtpDriver) SetPlainTextBody(body string) {
	// TODO implement
}
func (s *smtpDriver) SetAttachments(attachments []Attachment) {
	// TODO implement
}

func (s *smtpDriver) Send() {
	// TODO implement
}
