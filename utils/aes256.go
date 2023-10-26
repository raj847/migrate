package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"github.com/raj847/togrpc/config"
	"github.com/raj847/togrpc/models"
)

// Decrypt from base64 to decrypted string
func DecryptMerchantKey(cryptoText string, saltKey ...interface{}) (models.MerchantKey, error) {
	var result models.MerchantKey
	var merchKey models.MerchantKey
	keyText := config.SALT_KEY
	if len(saltKey) > 0 {
		keyText = saltKey[0].(string)
	}
	key := []byte(keyText)
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return result, err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return result, err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	unMarshall := json.Unmarshal(ciphertext, &merchKey)
	fmt.Println(unMarshall)
	result = merchKey
	return result, nil
}

// Encrypt string to base64 crypto using AES
func Encrypt(DataEncrypt interface{}, saltKey ...interface{}) (models.ResultEncryptMerchKey, error) {
	var result models.ResultEncryptMerchKey

	keyText := config.SALT_KEY
	if len(saltKey) > 0 {
		keyText = saltKey[0].(string)
	}

	key := []byte(keyText)

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(DataEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(reqBodyBytes.Bytes()))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return result, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], reqBodyBytes.Bytes())

	// convert to base64
	result.Result = base64.URLEncoding.EncodeToString(ciphertext)
	return result, nil
}

func Decrypt(cryptoText string, saltKey ...interface{}) (string, error) {
	var result string
	keyText := config.SALT_KEY
	if len(saltKey) > 0 {
		keyText = saltKey[0].(string)
	}

	key := []byte(keyText)
	cipherText, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		return result, err
	}

	if len(cipherText) < aes.BlockSize {
		return result, err
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(cipherText, cipherText)

	unMarshall := json.Unmarshal(cipherText, &result)
	fmt.Println(unMarshall)

	return result, nil
}
