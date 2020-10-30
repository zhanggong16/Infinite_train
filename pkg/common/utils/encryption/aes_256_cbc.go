package encryption

import (
	"fmt"
	"os/exec"
	"strings"
)

/*
  server should be install openssl
*/

const (
	key = "54460A61CFE68A53" // this key is openssl -K param
)

func Encrypt(value string, iv string) string {
	if value == "" {
		return ""
	}
	encrData, err := AES256CBCEncrypt(value, key, iv)
	if err != nil {
		fmt.Printf("encrypt fail value is %s, message:%+v", value, err)
		return value
	}
	return encrData
}

func Decrypt(value string, iv string) string {
	if value == "" {
		return ""
	}
	decrData, err := AES256CBCDecrypt(value, key, iv)
	if err != nil {
		fmt.Printf("decrypt fail value is %s, message:%+v", value, err)
		return value
	}
	return decrData
}

func AES256CBCEncrypt(origData string, key string, iv string) (string, error) {
	s := fmt.Sprintf("echo %s | openssl aes-256-cbc -K %s -base64 -iv %s -A", origData, key, iv)
	encrData, err := exec.Command("/bin/bash", "-c", s).Output()
	if err != nil {
		return "", err
	}
	encrDataStr := strings.TrimRight(string(encrData), "\n")
	return encrDataStr, nil
}

func AES256CBCDecrypt(encrData string, key string, iv string) (string, error) {
	s := fmt.Sprintf("echo %s | openssl aes-256-cbc -d -K %s -base64 -iv %s -A", encrData, key, iv)
	origData, err := exec.Command("/bin/bash", "-c", s).Output()
	if err != nil {
		return "", err
	}
	origDataStr := strings.TrimRight(string(origData), "\n")
	return origDataStr, nil
}