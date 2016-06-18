package common

import (
	logger "app/logger"
	"crypto/rand"
	"encoding/base64"
)

func CheckErr(err error, msg string) {
	if err != nil {
		logger.Error(msg, err)
		// log.Fatalln(msg, err)
	}
}

func GeneraToken8() string {
	token, err := GenerateRandomString(8)
	CheckErr(err, "GeneraToken16")
	return token
}

func GeneraToken16() string {
	token, err := GenerateRandomString(16)
	CheckErr(err, "GeneraToken16")
	return token
}

func GeneraToken32() string {
	token, err := GenerateRandomString(32)
	CheckErr(err, "GeneraToken32")
	return token
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
