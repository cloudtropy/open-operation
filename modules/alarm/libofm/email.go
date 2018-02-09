package libofm

import (
	"net/smtp"
	"strings"
)

const (
	HOST        = "smtp.cloutropy.com"
	SERVER_ADDR = "smtp.cloutropy.com:25"
	USER        = "console@cloutropy.com" //发送邮件的邮箱
	PASSWORD    = "Hello12345"            //发送邮件邮箱的密码
)

type Email struct {
	to      string "to"
	subject string "subject"
	msg     string "msg"
}

/*use unSSL to link mail server*/
type unencryptedAuth struct {
	smtp.Auth
}

func (a unencryptedAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	s := *server
	s.TLS = true
	return a.Auth.Start(&s)
}

func NewEmail(to, subject, msg string) *Email {
	return &Email{to: to, subject: subject, msg: msg}
}

func SendEmail(email *Email) error {
	auth := unencryptedAuth{
		smtp.PlainAuth(
			"",
			USER,
			PASSWORD,
			HOST,
		),
	}
	sendTo := strings.Split(email.to, ";")
	done := make(chan error, 1024)

	go func() {
		defer close(done)
		for _, v := range sendTo {

			str := strings.Replace("From: "+USER+"~To: "+v+"~Subject: "+email.subject+"~~", "~", "\r\n", -1) + email.msg

			err := smtp.SendMail(
				SERVER_ADDR,
				auth,
				USER,
				[]string{v},
				[]byte(str),
			)
			done <- err
		}
	}()

	for i := 0; i < len(sendTo); i++ {
		<-done
	}

	return nil
}
