# A Golang package for sending emails using SMTP, SparkPost, SendGrid and MailGun

## [Under Development]

## Features
- Multiple File attachments
- Multiple recipients
- Multiple CC
- Multiple BCC
- HTML content type support
- Plain Text content type support
- Easy integration with [smtp4dev](https://github.com/rnwood/smtp4dev/tree/master) testing server for development
- Multiple Drivers Support: SMTP, SparkPost, SendGrid and MailGun

## Install
Here is how to add it to `go.mod` in your project
```go
go get github.com/harranali/mailing
```
## Usage
Here is how to use it
```go
// Initiate the package with SMTP Driver
mailer := mailing.NewMailerWithSMTP(&mailing.SMTPConfig{
		Host:     "localhost", //the SMTP host
		Port:     25, // The Port
		Username: "", 
		Password: "",
		TLSConfig: tls.Config{
			ServerName:         "localhost",
			InsecureSkipVerify: true, // (use true for development only) true accepts any certificate presented by the server
		},
	})

// Set the Sender email address
mailer.SetFrom(mailing.EmailAddress{
        Name: "from name", 
        Address: "from@mail.com"
    })

// Set the Recipients email addresses
mailer.SetTo([]mailing.EmailAddress{
        {Name: "first to name", Address: "theFirstTo@mail.com"},
        {Name: "second to name", Address: "theSecondTo@mail.com"},
    })

// Set CC email addresses
mailer.SetCC([]mailing.EmailAddress{
        {Name: "cc name", Address: "cc1@mail.com"},
        {Name: "cc name", Address: "cc2@mail.com"},
    })

// Set BCC email addresses
mailer.SetBCC([]mailing.EmailAddress{
        {Name: "bcc name", Address: "bcc@mail.com"},
    })

// Set the subject
mailer.SetSubject("This is the subject")

// Set the body (pick one: either HTML or Plain Text)
mailer.SetHTMLBody("<h1>This is the email body</h1>")
// OR
mailer.SetPlainTextBody("This is the email body")

// Set the sttachments files
mailer.SetAttachments([]mailing.Attachment{
        {
            Name: "first file",
            Path: "./myfiles/first-file.jpg",
        },
        {
            Name: "second file",
            Path: "./myfiles/second-file.pdf",
        },
    })
        
// Send the email
err := mailer.Send()
if err != nil {
    panic(err.Error())
}
```