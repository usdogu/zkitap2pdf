package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

var iv = []byte("kxk101l0O1l0O1l0")

func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("invalid PKCS#7 padding (empty data)")
	}
	padding := int(data[len(data)-1])
	if padding == 0 || padding > len(data) {
		return nil, fmt.Errorf("invalid PKCS#7 padding")
	}
	return data[:len(data)-padding], nil
}

func Decrypt(encryptedBase64Reverse string) (string, error) {
	key := []byte("1l0O1l0O1l0O1l0O1l0O1l0O1l0O1l0O")

	encryptedBase64 := reverseString(encryptedBase64Reverse)

	encryptedData, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encryptedData))
	mode.CryptBlocks(decrypted, encryptedData)

	decrypted, err = pkcs7Unpad(decrypted)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func DecryptFile(keyPrefix string, contents []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(keyPrefix + "1l0O1l0O1l0O1l0O1l0O1l0O"))
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(contents))
	mode.CryptBlocks(decrypted, contents)

	decrypted, err = pkcs7Unpad(decrypted)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
