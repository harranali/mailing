// Copyright 2023 Harran Ali <harran.m@gmail.com>. All rights reserved.
// Use of this source code is governed by MIT-style
// license that can be found in the LICENSE file.

package mailing

import "net/mail"

type Driver interface {
	Send()
	SetFrom(from mail.Address)
	SetTo(toList []mail.Address)
	SetCC(ccList []mail.Address)
	SetBCC(bccList []mail.Address)
	SetSubject(subject string)
	SetHTMLBody(body string)
	SetPlainTextBody(body string)
	SetAttachments(attachments []Attachment)
}

type Mailer struct {
	driver     Driver
	sender     mail.Address
	receiver   mail.Address
	cc         []mail.Address
	bcc        []mail.Address
	subject    string
	htmlBody   string
	plainText  string
	attachment string
}

type EmailAddress struct {
	Name    string // the name can be empty
	Address string // ex: john@example.com
}
type Attachment struct {
	Name string // name of the file
	Path string // full path to the file
}

// Initiate the mailer with SMTP driver
func NewSMTPMailer(config *SMTPConfig) *Mailer {
	smtpDriver := initiateSMTP(config)
	return &Mailer{driver: smtpDriver}
}

// Initiate the mailer with SparkPost driver
func NewSparkPostMailer(config *SMTPConfig) *Mailer {
	// TODO implement
	return &Mailer{}
}

// Initiate the mailer with SendGrid driver
func NewSendGridMailer(config *SMTPConfig) *Mailer {
	// TODO implement
	return &Mailer{}
}

// Initiate the mailer with MailGun driver
func NewMailGunMailer(config *SMTPConfig) *Mailer {
	// TODO implement
	return &Mailer{}
}

// Sender of the email
func (m *Mailer) SetFrom(emailAddress EmailAddress) *Mailer {
	m.driver.SetFrom(mail.Address{Name: emailAddress.Name, Address: emailAddress.Address})
	return m
}

// List of receivers of the email
func (m *Mailer) SetTo(emailAddresses []EmailAddress) *Mailer {
	var AddressesList []mail.Address
	for _, v := range AddressesList {
		AddressesList = append(AddressesList, mail.Address{Name: v.Name, Address: v.Address})
	}

	m.driver.SetTo(AddressesList)
	return m
}

// List of cc of the email
func (m *Mailer) SetCC(emailAddresses []EmailAddress) *Mailer {
	var AddressesList []mail.Address
	for _, v := range AddressesList {
		AddressesList = append(AddressesList, mail.Address{Name: v.Name, Address: v.Address})
	}

	m.driver.SetCC(AddressesList)
	return m
}

// List of cc of the email
func (m *Mailer) SetBCC(emailAddresses []EmailAddress) *Mailer {
	var AddressesList []mail.Address
	for _, v := range AddressesList {
		AddressesList = append(AddressesList, mail.Address{Name: v.Name, Address: v.Address})
	}

	m.driver.SetBCC(AddressesList)
	return m
}

// Title of the email
func (m *Mailer) SetSubject(subject string) *Mailer {
	m.driver.SetSubject(subject)
	return m
}

// Set the body of the email in html format
// make sure to use only one version of the email body, either the html or the plain text
// to use the html, call the function SetHTMLBody(body string)
// and if you want to use the text, call the function SetPlainTextBody(body string)
func (m *Mailer) SetHTMLBody(body string) *Mailer {
	m.driver.SetHTMLBody(body)
	return m
}

// Set the body of the email in plain text format
// make sure to use only one version of the email body, either the html or the plain text
// to use the html, call the function SetHTMLBody(body string)
// and if you want to use the text, call the function SetPlainTextBody(body string)
func (m *Mailer) SetPlainTextBody(body string) *Mailer {
	m.driver.SetPlainTextBody(body)
	return m
}

// Add attachments to the email
func (m *Mailer) SetAttachments(attachments []Attachment) *Mailer {
	m.driver.SetAttachments(attachments)
	return m
}

// Send the email
func (m *Mailer) Send() (bool, error) {
	// TODO implement
	return false, nil
}
