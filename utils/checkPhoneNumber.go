package utils

import (
	"errors"
	"strconv"
	"strings"
)

// check by right number, potential case number
// 1. +628
// 2. 628
// 3. 08
// 4. 0823-1232-1323 (split karakter '-')
// dapatkan angka dari 822

func CheckPhoneNumber(phone string) (string, error) {
	if len(phone) < 8 {
		return "", errors.New("phone number length less than 8")
	}

	searchPhone := strings.ReplaceAll(phone, "-", "")

	if searchPhone[:4] == "+628" {
		return searchPhone[3:], nil
	}

	if searchPhone[:3] == "628" {
		return searchPhone[2:], nil
	}

	if searchPhone[:2] == "08" {
		return searchPhone[1:], nil
	}

	if _, err := strconv.Atoi(searchPhone); err != nil {
		return "", errors.New("phone number must be number")
	}

	return "", errors.New("phone number pattern not valid")
}
