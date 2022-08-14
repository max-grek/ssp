package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

func EncryptMD5(pwd []byte) string {
	hash := md5.Sum(pwd)
	return hex.EncodeToString(hash[:])
}

func EncryptAES(data, key []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("no source data")
	}
	c, err := aes.NewCipher([]byte(EncryptMD5(key)))
	if err != nil {
		return "", err
	}
	var gcm cipher.AEAD
	gcm, err = cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(gcm.Seal(nonce, nonce, []byte(data), nil)), nil
}

func DecryptAES(data, key []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("no source data")
	}
	dec, err := hex.DecodeString(string(data))
	if err != nil {
		return "", err
	}
	var (
		c   cipher.Block
		gcm cipher.AEAD
		out []byte
	)
	c, err = aes.NewCipher([]byte(EncryptMD5(key)))
	if err != nil {
		return "", err
	}
	gcm, err = cipher.NewGCM(c)
	if err != nil {
		return "", err
	}
	//Extract the nonce from the encrypted data
	nonce, ciphertext := dec[:gcm.NonceSize()], dec[gcm.NonceSize():]
	out, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
