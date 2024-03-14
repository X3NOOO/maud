package runners

import (
	"errors"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/X3NOOO/maud/types"
)

type Email struct {
	Host     string
	Port     int
	Email    string
	Password string

	auth smtp.Auth
}

func (e Email) Fire(sw types.Switch) error {
	if e.auth == nil {
		e.auth = smtp.PlainAuth("", e.Email, e.Password, e.Host)
	}
	if len(sw.Recipients) == 0 {
		return errors.New("no recipients")
	}

	var msg string

	msg += "From: " + e.Email + "\n"
	msg += "To: " + strings.Join(sw.Recipients, ", ") + "\n"
	msg += "Date: " + time.Now().Format(time.UnixDate) + "\n"
	msg += "Subject: " + sw.Subject + "\n"
	msg += "Importance: High" + "\n"
	msg += "X-Priority: 1 (Highest)" + "\n"
	msg += "X-MSMail-Priority: High" + "\n"
	msg += "\n"
	msg += sw.Content

	msg = strings.ReplaceAll(strings.ReplaceAll(msg, "\r\n", "\n"), "\n", "\r\n")

	return smtp.SendMail(e.Host+":"+strconv.Itoa(e.Port), e.auth, e.Email, sw.Recipients, []byte(msg))
}
