package app

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/rand"
	"time"
)

var secretkey = []byte("dsfewf64jwlj6so4difslkdj321")

const seqLength = 5

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateShortLink() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, seqLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateNewUserCookie() (string, string) {
	newID := make([]byte, 4)

	rand.Read(newID)
	encodedID := hex.EncodeToString(newID)

	h := hmac.New(sha256.New, secretkey)
	h.Write(newID)
	sign := h.Sum(nil)
	return encodedID, encodedID + hex.EncodeToString(sign)
}

func GetUserIDfromCookie(cookie string) (string, error) {
	data, err := hex.DecodeString(cookie)
	if err != nil {
		return "", err
	}
	id := data[:4]
	h := hmac.New(sha256.New, secretkey)
	h.Write(data[:4])
	sign := h.Sum(nil)

	if hmac.Equal(sign, data[4:]) {
		return hex.EncodeToString(id), nil
	} else {
		err := errors.New("sign check error")
		return "", err
	}
}
