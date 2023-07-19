package mailing

import "net/mail"

type messageBuilder struct {
	subject       string
	htmlBody      string
	plainTextBody string
	from          mail.Address
	toList        []mail.Address
	ccList        []mail.Address
	attachments   []Attachment
}

func newMessageBuilder() *messageBuilder {
	return &messageBuilder{
		subject:       "",
		htmlBody:      "",
		plainTextBody: "",
	}
}

func (m *messageBuilder) setSubject(subject string) *messageBuilder {
	m.subject = subject
	return m
}

func (m *messageBuilder) setHTMLBody(body string) *messageBuilder {
	m.htmlBody = body
	return m
}

func (m *messageBuilder) setPlainTextBody(body string) *messageBuilder {
	m.plainTextBody = body
	return m
}

func (m *messageBuilder) setFrom(from mail.Address) *messageBuilder {
	m.from = from
	return m
}

func (m *messageBuilder) setToList(toList []mail.Address) *messageBuilder {
	m.toList = toList
	return m
}

func (m *messageBuilder) setCCList(ccList []mail.Address) *messageBuilder {
	m.ccList = ccList
	return m
}

func (m *messageBuilder) setAttachments(attachments []Attachment) *messageBuilder {
	m.attachments = attachments
	return m
}

func (m *messageBuilder) build() string {
	// TODO implement
	return ""
}
