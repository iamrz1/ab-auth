package utils

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
)

func IsGenderValid(gender string) bool {
	genders := []string{"male", "female", "other"}
	if ExistsInSlice(genders, strings.ToLower(gender)) {
		return true
	}

	return false
}

func IsValidPhoneNumber(phoneNumber string) bool {
	if os.Getenv("ENV") == "dev" {
		log.Println(">============xxx NO PHONE NUMBER VALIDATION IN DEV BUILDS xxx============<")
		return true
	}
	reg := regexp.MustCompile(`^(01)[3-9][0-9]{8}$`)
	if !reg.MatchString(phoneNumber) {
		return false
	}

	return true
}

func ValidatePassword(password string) error {
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
