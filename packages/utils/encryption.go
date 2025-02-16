package utils

import "golang.org/x/crypto/bcrypt"

func IsValidPassword(hashedPw string, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(pw)) == nil
}

func HashString(word string) (string, error) {
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(word), 14)
	if err != nil {
		return "", err
	}

	return string(passwordByte), nil
}
