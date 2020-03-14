package utils

import (
	"regexp"
)

const (
	strRegEmail    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	strRegPhone    = "^(1[3|4|5|8])\\d{9}$"
	strRegUsername = "^[a-zA-Z]\\w{2,20}$"
)

var (
	regEmail    = regexp.MustCompile(strRegEmail)
	regPhone    = regexp.MustCompile(strRegPhone)
	regUsername = regexp.MustCompile(strRegUsername)
)

// IsEmail Check
func IsEmail(s string) bool {
	return regEmail.MatchString(s)
}

// IsPhone Check
func IsPhone(s string) bool {
	return regPhone.MatchString(s)
}

// // IsLetter Check
// func IsLetter(c byte) bool {
// 	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
// }

// IsUsername Check
func IsUsername(s string) bool {
	return regUsername.MatchString(s)
}

// CheckUsername as name/phone/email
func CheckUsername(s string) (name, phone, email string, ok bool) {
	if s == "" {
		return
	}
	if IsPhone(s) {
		phone = s
		ok = true
		return
	}
	if IsEmail(s) {
		email = s
		ok = true
		return
	}
	if IsUsername(s) {
		name = s
		ok = true
		return
	}
	return
}
