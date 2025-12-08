package util

import (
	"math/rand"
	"time"
)

const UPPER_CASE_LETTERS_AND_DIGITS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func StringGenerator(characterSet string, length int) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := make([]byte, length)
	for i := range password {
		password[i] = characterSet[rng.Intn(len(characterSet))]
	}

	return string(password)
}
