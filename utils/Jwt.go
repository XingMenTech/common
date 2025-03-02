package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

const (
	TokenSecret      = "game" //加密密钥
	TokenInvalidTime = 3      //小时
)

var tokenInfo string

type Claims struct {
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 校验时间
func (c *Claims) Valid() error {
	return c.StandardClaims.Valid()
}

var jwtSecret = []byte(TokenSecret)

// 生成token
func GenerateToken(uid int) string {
	token := "TK-" + fmt.Sprintf("%d", uid) + "-" + strings.ToUpper(RandomString(13))
	return token
}

//func GenerateToken(username string, uid int) (string, error) {
//now := time.Now()
//expireTime := now.Add(TokenInvalidTime * time.Hour)
//issuer := MD5(fmt.Sprintf("%d", now.Unix()))
//claims := Claims{
//	uid,
//	username,
//	jwt.StandardClaims{
//		ExpiresAt: expireTime.Unix(),
//		IssuedAt:  now.Unix(),
//		Issuer:    issuer,
//		NotBefore: now.Unix(),
//	},
//}
//tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
//
//token, err := tokenClaims.SignedString(jwtSecret)
//
//return token, err
//}

// 解析token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, errors.New("parse token error")
	}
	claims, ok := tokenClaims.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}
	// 验证是否过期
	if err1 := claims.Valid(); err1 != nil {
		return nil, err1
	}

	// 验证签发人
	issuer := MD5(fmt.Sprintf("%d", claims.StandardClaims.IssuedAt))
	if !claims.VerifyIssuer(issuer, true) {
		return nil, errors.New("invalid token issuer")
	}

	return claims, nil
}
