package client

import (
	"bytes"
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
	"strings"

	"git.championtek.com.tw/go/mlog"

	errors "git.championtek.com.tw/go/errors"
)

const (
	//SMTPServer Gmail SMTP Server
	SMTPServer = "smtp.gmail.com"
)

//Sender 定義寄送者資訊
type Sender struct {
	User     string
	Password string
}

//NewSender 回傳Sender資料
func NewSender(user, password string) Sender {
	return Sender{User: user, Password: password}
}

//SendMail 寄送信件，收件者為字串陣列型態
func (sender Sender) SendMail(recipients []string, subject, bodyMessage string) {
	msg := "From: " + sender.User + "\n" +
		"To: " + strings.Join(recipients, ",") + "\n" +
		"Subject: " + subject + "\n" + bodyMessage
	err := smtp.SendMail(
		SMTPServer+":587",
		smtp.PlainAuth("", sender.User, sender.Password, SMTPServer),
		sender.User, recipients, []byte(msg))
	if err != nil {
		errors.Wrap(err, "SendMail error")
		return
	}
	mlog.Info("Mail sent successfully!")
}

//WriteEmail 寫入信件資料
func (sender Sender) WriteEmail(recipients []string, contentType, subject, bodyMessage string) string {
	header := make(map[string]string)
	header["From"] = sender.User

	to := strings.Join(recipients, "")

	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", contentType)
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"

	message := ""

	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	var encodedMessage bytes.Buffer

	finaliMessage := quotedprintable.NewWriter(&encodedMessage)
	finaliMessage.Write([]byte(bodyMessage))
	finaliMessage.Close()

	message += "\r\n" + encodedMessage.String()

	return message
}

//WriteHTMLEmail 寫入HTML內容信件
func (sender *Sender) WriteHTMLEmail(recipients []string, subject, bodyMessage string) string {
	return sender.WriteEmail(recipients, "text/html", subject, bodyMessage)
}

//WritePlainEmail 寫入Plain Text內容信件
func (sender *Sender) WritePlainEmail(recipients []string, subject, bodyMessage string) string {
	return sender.WriteEmail(recipients, "text/plain", subject, bodyMessage)
}
