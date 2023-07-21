# A Golang package for sending emails using SMTP, SparkPost, SendGrid and MailGun

![Build Status](https://github.com/harranali/mailing/actions/workflows/build-main.yml/badge.svg)
![Test Status](https://github.com/harranali/mailing/actions/workflows/test-main.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/harranali/mailing)](https://goreportcard.com/report/github.com/harranali/mailing)
[![GoDoc](https://godoc.org/github.com/harranali/mailing?status.svg)](https://godoc.org/github.com/harranali/mailing)
[![Coverage Status](https://coveralls.io/repos/github/harranali/mailing/badge.svg?branch=main)](https://coveralls.io/github/harranali/mailing?branch=main&cache=false)

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
Here is how to add it to your project
```go
go get github.com/harranali/mailing
```
## Usage
Here is how to use it
```go
// Initiate the package with SMTP Driver
mailer := mailing.NewMailerWithSMTP(&mailing.SMTPConfig{
		Host:     "localhost", //the SMTP server host
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
        Name: "from name", // name can be set to empty string ("")
        Address: "from@mail.com"
    })

// Set the Recipients email addresses
mailer.SetTo([]mailing.EmailAddress{
        // name can be set to empty string ("")
        {Name: "first to name", Address: "theFirstTo@mail.com"},
        {Name: "second to name", Address: "theSecondTo@mail.com"},
    })

// Set CC email addresses
mailer.SetCC([]mailing.EmailAddress{
        // name can be set to empty string ("")
        {Name: "cc name", Address: "cc1@mail.com"},
        {Name: "cc name", Address: "cc2@mail.com"},
    })

// Set BCC email addresses
mailer.SetBCC([]mailing.EmailAddress{
        // name can be set to empty string ("")
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

## Testing you emails with smtp4dev SMTP Testing Server
While developing your app you might need to test your emails, for that a customized [docker-compose.yaml](https://github.com/harranali/mailing/tree/main/smtp-testing-server) from the SMTP testing server [smtp4dev](https://github.com/rnwood/smtp4dev/tree/master) is included.
#### Running the testing server
Copy the [docker-compose.yaml](https://github.com/harranali/mailing/blob/main/smtp-testing-server/docker-compose.yaml) to your pc, then start the container by running
```go
docker-compose up
```
#### The testing server configuration
Here is how to connect to the testing server
```go
mailer := mailing.NewMailerWithSMTP(&mailing.SMTPConfig{
		Host:     "localhost", //the SMTP server host
		Port:     25, // The Port
		Username: "", 
		Password: "",
		TLSConfig: tls.Config{
			ServerName:         "localhost",
			InsecureSkipVerify: true, // (use true for development only) true accepts any certificate presented by the server
		},
	})
```
#### Accessing the testing server UI
The testing server `UI` allows you to check the emails, here is a link to the ui of the testing server
```go
http://localhost:5000
```
![smtp4dev server ui](https://raw.githubusercontent.com/harranali/mailing/main/smtp-testing-server/screenshots/smtp4dev-server-ui.png "smtp server ui")