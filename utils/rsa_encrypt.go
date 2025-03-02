package utils

import (
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

// 生成32位MD5
func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func Sha256(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	sum := h.Sum(nil)
	sha := hex.EncodeToString(sum)
	return sha
}

func parsePrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key decode error")
	}
	pkcs1PrivateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("private key parse error")
	}
	return pkcs1PrivateKey.(*rsa.PrivateKey), nil
}

func parsePublicKey(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pkixPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := pkixPublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key type error")
	}
	return pub, nil
}

func PrivateKeySignAndBase64(privateKey []byte, data []byte) (string, error) {
	pkcs1PrivateKey, err := parsePrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	h := sha1.New()
	h.Write(data)
	hashed := h.Sum(nil)
	signPKCS1v15, err := rsa.SignPKCS1v15(nil, pkcs1PrivateKey, crypto.SHA1, hashed)
	if err != nil {
		return "", err
	}
	base64EncodingData := base64.StdEncoding.EncodeToString(signPKCS1v15)
	return base64EncodingData, nil
}

func PublicKeyEncryptAndBase64(src []byte, publicKey []byte) (string, error) {
	key, err := parsePublicKey(publicKey)
	if err != nil {
		return "", err
	}
	encryptPKCS1v15, err := rsa.EncryptPKCS1v15(rand.Reader, key, src)
	if err != nil {
		return "", err
	}
	base64EncodingData := base64.StdEncoding.EncodeToString(encryptPKCS1v15)
	return base64EncodingData, nil
}
