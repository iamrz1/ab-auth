package utils

import (
	"log"
	"os"
	"regexp"
	"strings"
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

