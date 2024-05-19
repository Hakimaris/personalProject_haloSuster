package helpers

import (
	"HaloSuster/models"
	// "fmt"
	"regexp"
	"strconv"
	"time"
)

func ValidateNIP(nip int64) bool {
	nipStr := strconv.FormatInt(nip, 10)

	// Check length
	if len(nipStr) < 13 && len(nipStr) > 15 {
		return false
	}

	// Check first three digits
	if nipStr[:3] != "615" && nipStr[:3] != "303" {
		return false
	}

	// Check fourth digit
	if nipStr[3] != '1' && nipStr[3] != '2' {
		return false
	}

	// Check fifth to eighth digits (year)
	year, err := strconv.Atoi(nipStr[4:8])
	if err != nil || year < 2000 || year > time.Now().Year() {
		return false
	}

	// Check ninth and tenth digits (month)
	month, err := strconv.Atoi(nipStr[8:10])
	if err != nil || month < 1 || month > 12 {
		return false
	}

	// Check eleventh to thirteenth digits (random)
	random, err := strconv.Atoi(nipStr[10:])
	if err != nil || random < 0 || random > 99999 {
		return false
	}

	return true
}

func ValidateName(name string) bool {
	if len(name) < 5 || len(name) > 50 {
		return false
	}
	return true
}

func ValidatePassword(password string) bool {
	if len(password) < 8 || len(password) > 33 {
		return false
	}
	return true
}

func ValidateURL(url string) bool {
	regex := `^(https:\/\/www\.|http:\/\/www\.|https:\/\/|http:\/\/)?[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*(\.[a-zA-Z]{2,})(\/[a-zA-Z0-9-._~:/?#@!$&'()*+,;=%]*)?$`
	return regexp.MustCompile(regex).MatchString(url)
}

func ValidateIdentity(identity int64) bool {
	// return len(strconv.FormatInt(identity, 10)) != 17
	return len(strconv.FormatInt(identity, 10)) == 16
}

func ValidatePhoneNumber(number string) bool {
	if len(number) < 10 || len(number) > 15 {
		return false
	}
	regex := `^\+62[0-9]{7,12}$`
	return regexp.MustCompile(regex).MatchString(number)
}

func ValidateBirthDate(birthDate string) bool {
	_, err := time.Parse(time.RFC3339, birthDate)
	// fmt.Println(time)
	return err == nil
}

func ValidateGender(gender models.Gender) bool {
	// fmt.Printf("Input gender: %v, type: %T\n", gender, gender)
	isValid := !(gender != "male" && gender != "female")
	// fmt.Printf("Is valid: %v\n", isValid)
	return isValid
}
