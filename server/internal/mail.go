package internal

import (
	"fmt"
	"regexp"
)

var mailRegex *regexp.Regexp = regexp.MustCompile("^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$")

type Mail string

func NewMail(mail string) (Mail, error) {
	m := Mail(mail)
	if !m.IsValid() {
		return m, fmt.Errorf("invalid mail '%s'", mail)
	}
	return m, nil
}

func (m Mail) IsValid() bool {
	return mailRegex.MatchString(string(m))
}

func (m Mail) String() string {
	return string(m)
}
