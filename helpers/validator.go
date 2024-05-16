package helpers

import (
	"strconv"
	"time"
)

func ValidateNIP(nip int64) bool {
	nipStr := strconv.FormatInt(nip, 10)

	// Check length
	if len(nipStr) != 13 {
		return false
	}

	// Check first three digits
	if nipStr[:3] != "615" && nipStr[:3] != "303"{
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
	if err != nil || random < 0 || random > 999 {
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

func ValidateRoleRequest(role string) bool {
	if role != "nurse" && role != "it" {
		return false
	}
	return true
}

func ValidateCreatedAtRequest(created string) bool {
	if created != "asc" && created != "desc" {
		return false
	}
	return true
}

