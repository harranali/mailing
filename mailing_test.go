package mailing

import (
	"crypto/tls"
	"testing"
)

func TestMailingParamSetters(t *testing.T) {
	mailer := NewMailerWithSMTP(&SMTPConfig{
		Host:     "localhost",
		Port:     25,
		Username: "",
		Password: "",
		TLSConfig: tls.Config{
			ServerName:         "localhost",
			InsecureSkipVerify: true,
		},
	})
	mailer.
		SetFrom(EmailAddress{Name: "from name", Address: "from@mail.com"}).
		SetTo([]EmailAddress{
			{Name: "to name", Address: "theFirstTo@mail.com"},
			{Name: "to name", Address: "theSecondTo@mail.com"},
		}).
		SetCC([]EmailAddress{
			{Name: "cc name", Address: "cc1@mail.com"},
			{Name: "cc name", Address: "cc2@mail.com"},
		}).
		SetBCC([]EmailAddress{
			{Name: "bcc name", Address: "bcc@mail.com"},
		}).
		SetSubject("This is the subject").
		SetHTMLBody("this is the body").
		SetAttachments([]Attachment{
			{
				Name: "attachment name1",
				Path: "./testingdata/attachment1.md",
			},
			{
				Name: "attachment name2",
				Path: "./testingdata/attachment2.md",
			},
		})

	dr := mailer.driver.(*smtpDriver)
	if !(dr.from.Name == "from name" && dr.from.Address == "from@mail.com") {
		panic("failed testing mailing parameters setting")
	}
	if len(dr.toList) != 2 {
		panic("failed testing mailing parameters setting")
	}
	if len(dr.ccList) != 2 {
		panic("failed testing mailing parameters setting")
	}
	if len(dr.bccList) != 1 {
		panic("failed testing mailing parameters setting")
	}
	if dr.subject != "This is the subject" {
		panic("failed testing mailing parameters setting")
	}
	if dr.htmlBody != "this is the body" {
		panic("failed testing mailing parameters setting")
	}
	if len(dr.attachments) != 2 {
		panic("failed testing mailing parameters setting")
	}
}
