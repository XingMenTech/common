package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
)

type AES struct {
	Key []byte
	Iv  []byte
}

func NewAES(key []byte, iv ...[]byte) *AES {
	var iv_ []byte
	if len(iv) > 0 && len(iv[0]) > 0 {
		iv_ = iv[0]
	} else {
		iv_ = key
	}
	return &AES{Key: key, Iv: iv_}
}

// Encrypt
func (a *AES) Encrypt(origData []byte) (string, error) {

	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData = pKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, a.Iv[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	// hex encode
	encode := hex.EncodeToString(crypted)
	return encode, nil
}

// Decrypt
func (a *AES) Decrypt(crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.Key)
	if err != nil {
		return nil, err
	}
	if len(crypted) <= 0 {
		return nil, errors.New("not enough length")
	}
	// hex decode
	if crypted, err = hex.DecodeString(string(crypted)); err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	if len(crypted)%blockSize != 0 {
		return nil, errors.New("input not full block!")
	}

	blockMode := cipher.NewCBCDecrypter(block, a.Iv[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData, err = pKCS7UnPadding(origData)
	return origData, err
}

func pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func pKCS7UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unpadding := int(origData[length-1])
	if length < unpadding {
		return origData, errors.New("slice bounds out of range")
	}
	return origData[:(length - unpadding)], nil
}
