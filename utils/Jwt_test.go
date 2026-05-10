package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func TestGenerateAndValidateAccessToken(t *testing.T) {
	// 生成访问令牌
	accessToken, err := GenerateAccessToken(123, "testuser", 1, "device123", "192.168.1.1")
	if err != nil {
		t.Fatalf("Failed to generate access token: %v", err)
	}

	// 验证访问令牌
	claims, err := ValidateAccessToken(accessToken)
	if err != nil {
		t.Fatalf("Failed to validate access token: %v", err)
	}

	// 验证声明内容
	if claims.Uid != 123 {
		t.Errorf("Expected UID 123, got %d", claims.Uid)
	}
	if claims.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", claims.Username)
	}
	if claims.UserType != 1 {
		t.Errorf("Expected user type 1, got %d", claims.UserType)
	}
	if claims.DeviceId != "device123" {
		t.Errorf("Expected device ID 'device123', got '%s'", claims.DeviceId)
	}
	if claims.ClientIp != "192.168.1.1" {
		t.Errorf("Expected client IP '192.168.1.1', got '%s'", claims.ClientIp)
	}
}

func TestGenerateTokenPair(t *testing.T) {
	accessToken, refreshToken, err := GenerateTokenPair(456, "anotheruser", 2, "device456", "10.0.0.1")
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	if accessToken == "" {
		t.Error("Access token should not be empty")
	}
	if refreshToken == "" {
		t.Error("Refresh token should not be empty")
	}

	// 验证访问令牌
	accessClaims, err := ValidateAccessToken(accessToken)
	if err != nil {
		t.Fatalf("Failed to validate access token: %v", err)
	}
	if accessClaims.Uid != 456 {
		t.Errorf("Expected UID 456, got %d", accessClaims.Uid)
	}

	// 验证刷新令牌
	refreshClaims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		t.Fatalf("Failed to validate refresh token: %v", err)
	}
	if refreshClaims.Uid != 456 {
		t.Errorf("Expected UID 456 in refresh token, got %d", refreshClaims.Uid)
	}
}

func TestExpiredToken(t *testing.T) {
	// 创建一个立即过期的令牌用于测试
	now := time.Now()
	expireTime := now.Add(-1 * time.Hour) // 设置为过去的时间使令牌立即过期

	claims := Claims{
		Uid:      789,
		Username: "expiringuser",
		UserType: 1,
		DeviceId: "device789",
		ClientIp: "192.168.1.100",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			Issuer:    "common-jwt-service",
			Subject:   "user_789",
			Id:        RandomString(32),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		t.Fatalf("Failed to generate expired token: %v", err)
	}

	// 尝试验证已过期的令牌
	_, err = ValidateAccessToken(accessToken)
	if err == nil {
		t.Error("Expected error for expired token, but got none")
	}
}

func TestInvalidToken(t *testing.T) {
	// 尝试验证无效令牌
	_, err := ValidateAccessToken("invalid.token.here")
	if err == nil {
		t.Error("Expected error for invalid token, but got none")
	}
}

func TestRefreshTokenFlow(t *testing.T) {
	// 生成初始令牌对
	accessToken, refreshToken, err := GenerateTokenPair(101, "refreshuser", 1, "devicerefresh", "192.168.1.50")
	if err != nil {
		t.Fatalf("Failed to generate initial token pair: %v", err)
	}

	// 验证初始访问令牌
	initialClaims, err := ValidateAccessToken(accessToken)
	if err != nil {
		t.Fatalf("Failed to validate initial access token: %v", err)
	}
	if initialClaims.Uid != 101 {
		t.Errorf("Expected UID 101, got %d", initialClaims.Uid)
	}

	// 使用刷新令牌获取新的令牌对
	newAccessToken, newRefreshToken, err := RefreshAccessToken(refreshToken)
	if err != nil {
		t.Fatalf("Failed to refresh access token: %v", err)
	}

	// 验证新的访问令牌
	newClaims, err := ValidateAccessToken(newAccessToken)
	if err != nil {
		t.Fatalf("Failed to validate new access token: %v", err)
	}
	if newClaims.Uid != 101 {
		t.Errorf("Expected UID 101 in new token, got %d", newClaims.Uid)
	}

	// 验证新的刷新令牌
	newRefreshClaims, err := ValidateRefreshToken(newRefreshToken)
	if err != nil {
		t.Fatalf("Failed to validate new refresh token: %v", err)
	}
	if newRefreshClaims.Uid != 101 {
		t.Errorf("Expected UID 101 in new refresh token, got %d", newRefreshClaims.Uid)
	}
}

func TestGetUserInfoFromToken(t *testing.T) {
	accessToken, _, err := GenerateTokenPair(202, "infouser", 3, "deviceinfo", "192.168.1.200")
	if err != nil {
		t.Fatalf("Failed to generate token pair: %v", err)
	}

	uid, err := GetUidFromToken(accessToken)
	if err != nil {
		t.Fatalf("Failed to get UID from token: %v", err)
	}
	if uid != 202 {
		t.Errorf("Expected UID 202, got %d", uid)
	}

	username, err := GetUsernameFromToken(accessToken)
	if err != nil {
		t.Fatalf("Failed to get username from token: %v", err)
	}
	if username != "infouser" {
		t.Errorf("Expected username 'infouser', got '%s'", username)
	}
}
