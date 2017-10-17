package gosmail

import "testing"

func TestMail(t *testing.T) {
	c := NewClient("smtp.example.com", ":587", "arion@example.com", "heljuurtupmebfch")

	m := new(Mail)
	m.From.Address = "arion@example.com"
	m.From.Name = "this is mail for test"
	m.To = []string{"kate@example.com"}
	m.Subject = "this is love letter"
	m.Body = "hi"
	m.HTMLBody = "<h1>title</h1>"
	m.Cc = []string{"marry@example.com"}
	attach := &Attachment{
		FileName:    "a.png",
		ContentType: "image/png",
		Data:        make([]byte, 5),
	}
	attach2 := &Attachment{
		FileName:    "a.txt",
		ContentType: "text/plain",
		Data:        make([]byte, 10),
	}
	attach3 := &Attachment{
		IsInline:    true,
		FileName:    "b.png",
		ContentType: "image/png",
		Data:        []byte("123123"),
	}

	m.Attachments = append(m.Attachments, attach, attach2, attach3)
	c.Send(m)
}
