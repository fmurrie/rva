package helper

import (
	"crypto/aes"
	"crypto/cipher"
	cryptoRand "crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var singletonSecurityHelper *SecurityHelper
var onceSecurityHelper sync.Once

type SecurityHelper struct {
	securityKey string
}

func GetSecurityHelper() *SecurityHelper {
	onceSecurityHelper.Do(func() {
		singletonSecurityHelper = &SecurityHelper{
			securityKey: hex.EncodeToString(make([]byte, 32)),
		}
	})
	return singletonSecurityHelper
}

func (pointer SecurityHelper) Encrypt(stringToEncrypt string) string {
	stringEncrypted := ""

	randomOption := func(options ...string) string {
		rand.Seed(time.Now().UnixNano())
		return options[rand.Intn(len(options))]
	}

	for _, value := range stringToEncrypt {
		asciiValue := fmt.Sprint(int(value))
		for _, val := range asciiValue {
			switch val {
			case '0':
				stringEncrypted = fmt.Sprint(stringEncrypted, randomOption("В", "X", "K"))
			case '1':
				stringEncrypted = fmt.Sprint(stringEncrypted, randomOption("А", "T", "Β"))
			case '2':
				stringEncrypted = fmt.Sprint(stringEncrypted, randomOption("E", "H", "Ο"))
			case '3':
				stringEncrypted = fmt.Sprint(stringEncrypted, randomOption("О", "M", "Ε"))
			case '4':
				stringEncrypted = fmt.Sprint(stringEncrypted, randomOption("C", "Κ", "Ρ"))
			case '5':
				stringEncrypted = fmt.Sprint(stringEncrypted, randomOption("Р", "Χ", "С"))
			case '6':
				stringEncrypted = fmt.Sprint(stringEncrypted, randomOption("Н", "Е", "Μ"))
			case '7':
				stringEncrypted = fmt.Sprint(stringEncrypted, randomOption("М", "O", "Η"))
			case '8':
				stringEncrypted = fmt.Sprint(stringEncrypted, randomOption("Х", "B", "Τ"))
			case '9':
				stringEncrypted = fmt.Sprint(stringEncrypted, randomOption("Т", "Α", "P"))
			}
		}
		stringEncrypted = fmt.Sprint(stringEncrypted, "A")
	}

	key, _ := hex.DecodeString(pointer.securityKey)
	plaintext := []byte(stringEncrypted)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(cryptoRand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	return fmt.Sprintf("%x", ciphertext)
}

func (pointer SecurityHelper) Decrypt(encryptedString string) string {
	key, _ := hex.DecodeString(pointer.securityKey)
	enc, _ := hex.DecodeString(encryptedString)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	encryptedString = fmt.Sprintf("%s", plaintext)

	decryptedString := ""
	asciiString := ""
	for _, value := range encryptedString {
		switch string(value) {
		case "A":
			if asciiString != "" {
				ascii, _ := strconv.Atoi(asciiString)
				decryptedString = decryptedString + string(ascii)
				asciiString = ""
			}
		case "В", "X", "K":
			asciiString = asciiString + "0"
		case "А", "T", "Β":
			asciiString = asciiString + "1"
		case "E", "H", "Ο":
			asciiString = asciiString + "2"
		case "О", "M", "Ε":
			asciiString = asciiString + "3"
		case "C", "Κ", "Ρ":
			asciiString = asciiString + "4"
		case "Р", "Χ", "С":
			asciiString = asciiString + "5"
		case "Н", "Е", "Μ":
			asciiString = asciiString + "6"
		case "М", "O", "Η":
			asciiString = asciiString + "7"
		case "Х", "B", "Τ":
			asciiString = asciiString + "8"
		case "Т", "Α", "P":
			asciiString = asciiString + "9"
		}
	}
	return decryptedString
}
