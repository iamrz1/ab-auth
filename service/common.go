package service

import (
	"fmt"
	"unicode"
)

const (
	PasswordPattern         = "^([a-zA-z0-9!@#%*_=+/-]*)$"
	InvalidCharErrorMessage = "Password contains invalid characters"
	SpecialCharErrorMessage = "Must contain at least one special character"
	CharLenErrorMessage     = "Must be at least 8 characters long"
)

func validatePassword(password string) error {
	var eightOrMore, special, invalid bool
	var err error
	characters := 0
	for _, c := range password {
		switch {
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
			break
		case !(unicode.IsLetter(c) || unicode.IsNumber(c)):
			invalid = true
		}

		characters++
	}
	eightOrMore = characters >= 8

	if !special {
		err = fmt.Errorf("%s", SpecialCharErrorMessage)
	}
	if !eightOrMore {
		err = fmt.Errorf("%s", CharLenErrorMessage)
	}
	if invalid {
		err = fmt.Errorf("%s", InvalidCharErrorMessage)
	}

	return err
}
