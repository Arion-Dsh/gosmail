// mgosmail is a text/html/attachments mail Sender for human.

package gosmail

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"
)

//Mail a mail for send
type Mail struct {
	From        mail.Address // from mail, must be the same with auth
	To          []string     // the mail send to
	Cc          []string     // the cc may be empty.
	Bcc         []string     // the same with cc
	ReplyTo     string
	Subject     string        // the mail Subject
	Body        string        // the text mail body, may be empty
	HTMLBody    string        // the html mail body, may be empty
	Attachments []*Attachment //the Attachments include inline attachment
}

func (m *Mail) writeHeaHer(buf *bytes.Buffer) {
	// from
	buf.WriteString("From: " + m.From.String() + "\r\n")

	//date
	t := time.Now()
	buf.WriteString("Date: " + t.Format(time.RFC1123Z) + "\r\n")

	// to
	m.To = append(m.To, m.Cc...)
	m.To = append(m.To, m.Bcc...)
	buf.WriteString("To: " + strings.Join(m.To, ",") + "\r\n")

	//cc
	if len(m.Cc) > 0 {
		buf.WriteString("CC: " + strings.Join(m.Cc, ",") + "\r\n")

	}
	//bcc
	if len(m.Bcc) > 0 {
		buf.WriteString("Bcc: " + strings.Join(m.Bcc, ",") + "\r\n")

	}

	//Subject
	buf.WriteString("Subject: " + m.Subject + "\r\n")

	//Reply
	if len(m.ReplyTo) > 0 {
		buf.WriteString("Reply-To: " + m.ReplyTo + "\r\n")

	}

}

func (m *Mail) writerBody(buf *bytes.Buffer, mixed *multipart.Writer) {

	part := multipart.NewWriter(buf)
	defer part.Close()

	alt := make(textproto.MIMEHeader, 1)
	c := fmt.Sprintf("multipart/alternative; \r\n\tboundary=\"%s\"", part.Boundary())
	alt.Add("Content-Type", c)
	mixed.CreatePart(alt)

	write := func(c, t string) {
		ct := fmt.Sprintf("%s; charset=UTF-8", c)
		h := textproto.MIMEHeader{"Content-Type": []string{ct}}
		w, _ := part.CreatePart(h)
		w.Write([]byte(t))
	}
	write("text/plain", m.Body)
	write("text/html", m.HTMLBody)
}

func (m *Mail) toBytes() []byte {
	buf := bytes.NewBuffer(nil)

	//writeHeaHer
	m.writeHeaHer(buf)

	// Start our multipart/mixed part
	mixed := multipart.NewWriter(buf)
	defer mixed.Close()

	fmt.Fprintf(buf, "Content-Type: multipart/mixed;\r\n\tboundary=\"%s\"; charset=UTF-8\r\n\r\n", mixed.Boundary())

	m.writerBody(buf, mixed)

	// writer Attachments
	for _, attachment := range m.Attachments {
		w, _ := mixed.CreatePart(attachment.getMIMEHeader())
		b := make([]byte, base64.StdEncoding.EncodedLen(len(attachment.Data)))
		base64.StdEncoding.Encode(b, attachment.Data)
		w.Write(b)

	}
	return buf.Bytes()
}

//Attachment  the attachment
type Attachment struct {
	IsInline    bool
	FileName    string
	ContentType string
	Data        []byte
}

func (a *Attachment) getMIMEHeader() textproto.MIMEHeader {

	header := make(textproto.MIMEHeader, 1)

	header.Add("Content-Transfer-Encoding", "base64")
	header.Add("Content-Type", a.ContentType)

	if a.IsInline {
		disp := fmt.Sprintf("inline;\n\tfilename=%s", a.FileName)
		header.Add("Content-Disposition", disp)

	} else {
		disp := fmt.Sprintf("attachment; filename=%s", a.FileName)
		cid := fmt.Sprintf("<%s>", a.FileName)
		header.Add("Content-Disposition", disp)
		header.Add("cid", cid)
	}

	return header
}

//Client the smtp client
type Client struct {
	auth     smtp.Auth
	hostname string
	port     string
}

//Send send mail
func (c *Client) Send(mail *Mail) {

	err := smtp.SendMail(c.hostname+c.port, c.auth, mail.From.Address, mail.To, mail.toBytes())
	if err != nil {
		panic(err)
	}
}

//NewClient return a new smtp  client with auth plain
func NewClient(hostname, port, user, pwd string) *Client {
	c := new(Client)
	c.hostname = hostname
	c.port = port
	c.auth = smtp.PlainAuth("", user, pwd, hostname)

	return c
}
