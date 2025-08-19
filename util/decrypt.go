package util

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

var iv = []byte("kxk101l0O1l0O1l0")

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
