package mail

import (
	"cathub.me/go-web-examples/pkg/setting"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/rs/zerolog/log"
	"net/smtp"
)

func SendEmail(subject, text string, to ...string) error {
	conf := setting.Email

	mail := email.NewEmail()
	mail.From = conf.Username
	mail.To = to
	mail.Subject = subject
	mail.Text = []byte(text)
	err := mail.Send(fmt.Sprintf("%s:%d", conf.Host, conf.Port), smtp.PlainAuth("", conf.Username, conf.Password, conf.Host))
	if err != nil {
		log.Err(err).Msgf("Send email failed, to: %v", to)
		return err
	}
	return nil
}
