package mailing

import (
	"net/mail"
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {
	m := newMessageBuilder()
	m.setFrom(mail.Address{
		Name:    "from test name",
		Address: "from@mail.com",
	})
	m.setToList([]mail.Address{
		{
			Name:    "to test name1",
			Address: "to1@mail.com",
		},
		{
			Name:    "to test name2",
			Address: "to2@mail.com",
		},
	})
	m.setCCList([]mail.Address{
		{
			Name:    "tcc test name1",
			Address: "cc1@mail.com",
		},
		{
			Name:    "cc test name2",
			Address: "cc2@mail.com",
		},
	})
	m.setSubject("the subject")
	m.setHTMLBody("this is html body")
	m.setAttachments([]Attachment{
		{
			Name: "attachment name1",
			Path: "./testingdata/attachment1.md",
		},
		{
			Name: "attachment name2",
			Path: "./testingdata/attachment2.md",
		},
	})
	message := string(m.build())
	if !strings.Contains(message, `From: "from test name" <from@mail.com>`) {
		t.Error("Failed test build")
	}
	if !strings.Contains(message, `To: "to test name1" <to1@mail.com>;"to test name2" <to2@mail.com>`) {
		t.Error("Failed test build")
	}
	if !strings.Contains(message, `CC: "tcc test name1" <cc1@mail.com>;"cc test name2" <cc2@mail.com>`) {
		t.Error("Failed test build")
	}
	if !strings.Contains(message, `Subject: the subject`) {
		t.Error("Failed test build")
	}
	if !strings.Contains(message, `Content-Type: text/html; charset="UTF-8"`) {
		t.Error("Failed test build")
	}
	if !strings.Contains(message, `this is html body`) {
		t.Error("Failed test build")
	}
	if !strings.Contains(message, `Content-Disposition: attachment; filename="attachment name1"`) {
		t.Error("Failed test build")
	}
	if !strings.Contains(message, `Content-Disposition: attachment; filename="attachment name1"`) {
		t.Error("Failed test build")
	}
	if !strings.Contains(message, `dGhpcyBpcyBhIHRlc3QgZmlsZSBmb3IgZW1haWwgYXR0YWNobWVudCAx`) {
		t.Error("Failed test build")
	}
	if !strings.Contains(message, `Content-Disposition: attachment; filename="attachment name2"`) {
		t.Error("Failed test build")
	}
	if !strings.Contains(message, `dGhpcyBpcyBhIHRlc3QgZmlsZSBmb3IgZW1haWwgYXR0YWNobWVudCAy`) {
		t.Error("Failed test build")
	}
	m.setHTMLBody("")
	m.setPlainTextBody("this is plain text body")
	message = string(m.build())
	if !strings.Contains(message, `Content-Type: text/plain; charset="UTF-8"`) {
		t.Error("Failed test build")
	}
	if !strings.Contains(message, `this is plain text body`) {
		t.Error("Failed test build")
	}
}
