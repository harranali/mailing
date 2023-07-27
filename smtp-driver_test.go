package mailing

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestSMTPDriverSend(t *testing.T) {
	sDriver := initiateSMTP(&SMTPConfig{
		Host:     "localhost",
		Port:     25,
		Username: "",
		Password: "",
		TLSConfig: tls.Config{
			ServerName:         "localhost",
			InsecureSkipVerify: true,
		},
	})
	tmpFilePath := filepath.Join(t.TempDir(), uuid.NewString())
	sDriver.initiateSend = func(from string, rcpts []string, message []byte, d Driver) error {
		file, err := os.Create(tmpFilePath)
		if err != nil {
			t.Error("faild test send")
		}
		file.Write(message)
		file.Close()
		return nil
	}

	sDriver.SetFrom(mail.Address{
		Name:    "test from name",
		Address: "from@mail.com",
	})
	sDriver.SetTo([]mail.Address{
		{Name: "test from name1", Address: "from1@mail.com"},
		{Name: "test from name2", Address: "from2@mail.com"},
	})
	sDriver.SetCC([]mail.Address{
		{Name: "test cc name1", Address: "cc1@mail.com"},
		{Name: "test cc name2", Address: "cc2@mail.com"},
	})
	sDriver.SetBCC([]mail.Address{
		{Name: "test bcc name1", Address: "bcc1@mail.com"},
		{Name: "test bcc name2", Address: "bcc2@mail.com"},
	})
	sDriver.SetSubject("this is the subject")
	sDriver.SetPlainTextBody("this is plain text body")
	sDriver.SetHTMLBody("this is html body")
	sDriver.SetAttachments([]Attachment{
		{
			Name: "attachment name1",
			Path: "./testingdata/attachment1.md",
		},
		{
			Name: "attachment name2",
			Path: "./testingdata/attachment2.md",
		},
	})
	err := sDriver.Send()
	if err != nil {
		t.Error("failed testing send")
	}
	file, _ := os.Open(tmpFilePath)
	mBytes, _ := ioutil.ReadAll(file)
	file.Close()
	os.Truncate(tmpFilePath, 0)
	m := string(mBytes)

	if !strings.Contains(m, `From: "test from name" <from@mail.com>`) {
		t.Error("Failed test send")
	}
	if !strings.Contains(m, `To: "test from name1" <from1@mail.com>;"test from name2" <from2@mail.com>`) {
		t.Error("Failed test send")
	}
	if !strings.Contains(m, `Cc: "test cc name1" <cc1@mail.com>;"test cc name2" <cc2@mail.com>`) {
		t.Error("Failed test send")
	}
	if !strings.Contains(m, `Subject: this is the subject`) {
		t.Error("Failed test send")
	}
	if !strings.Contains(m, `Content-Type: text/html; charset="UTF-8"`) {
		t.Error("Failed test send")
	}
	if !strings.Contains(m, `this is html body`) {
		t.Error("Failed test send")
	}
	if !strings.Contains(m, `Content-Disposition: attachment; filename="attachment name1"`) {
		t.Error("Failed test send")
	}
	if !strings.Contains(m, `dGhpcyBpcyBhIHRlc3QgZmlsZSBmb3IgZW1haWwgYXR0YWNobWVudCAx`) {
		t.Error("Failed test send")
	}
	if !strings.Contains(m, `Content-Disposition: attachment; filename="attachment name2"`) {
		t.Error("Failed test send")
	}
	if !strings.Contains(m, `dGhpcyBpcyBhIHRlc3QgZmlsZSBmb3IgZW1haWwgYXR0YWNobWVudCAy`) {
		t.Error("Failed test send")
	}

	sDriver.SetHTMLBody("")
	sDriver.SetPlainTextBody("this is plain text body")
	err = sDriver.Send()
	if err != nil {
		t.Error("failed testing send")
	}
	file, _ = os.Open(tmpFilePath)
	mBytes, _ = ioutil.ReadAll(file)
	m = string(mBytes)
	if !strings.Contains(m, `Content-Type: text/plain; charset="UTF-8"`) {
		t.Error("Failed test send")
	}
	if !strings.Contains(m, `this is plain text body`) {
		t.Error("Failed test send")
	}

	sDriver.initiateSend = func(from string, rcpts []string, message []byte, d Driver) error {
		return errors.New("this is a test error")
	}
	err = sDriver.Send()
	if err == nil {
		t.Error("failed testing send")
	}
}
