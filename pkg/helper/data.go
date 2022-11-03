package helper

import "net/mail"

func IsEmpty(param string) bool {
	return param == ""
}

func ErrorIsNil(err error) bool {
	return err != nil
}

func NotMatch(a, b string) bool {
	return a != b
}

func CheckPatternEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
