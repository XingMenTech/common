package utils

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ContextKey 定义上下文键
type ContextKey string

const (
	UserClaimsContextKey ContextKey = "user_claims"
)

// JWTMiddleware JWT认证中间件
type JWTMiddleware struct {
	// 可选：自定义验证函数
	CustomValidator func(*Claims) error
	// 可选：排除的路径列表
	ExcludedPaths []string
}

// NewJWTMiddleware 创建新的JWT中间件
func NewJWTMiddleware() *JWTMiddleware {
	return &JWTMiddleware{
		ExcludedPaths: []string{},
	}
}

// GinMiddleware Gin框架兼容的中间件
func (m *JWTMiddleware) GinMiddleware() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		// 检查是否在排除路径中
		path := c.Request.URL.Path
		for _, excluded := range m.ExcludedPaths {
			if path == excluded || strings.HasPrefix(path, excluded+"/") {
				c.Next()
				return
			}
		}

		// 获取Authorization头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "missing authorization header",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "invalid authorization format",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 验证令牌
		claims, err := ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "invalid or expired token: " + err.Error(),
			})
			c.Abort()
			return
		}

		// 如果有自定义验证器，则执行
		if m.CustomValidator != nil {
			if err := m.CustomValidator(claims); err != nil {
				c.JSON(http.StatusForbidden, map[string]interface{}{
					"code":    403,
					"message": "access denied: " + err.Error(),
				})
				c.Abort()
				return
			}
		}

		// 将用户信息存入上下文
		c.Set(string(UserClaimsContextKey), claims)
		c.Next()
	}
}

// GetUserFromGinContext 从Gin上下文中获取用户声明
func GetUserFromGinContext(c *gin.Context) (*Claims, bool) {
	claims, exists := c.Get(string(UserClaimsContextKey))
	if !exists {
		return nil, false
	}
	userClaims, ok := claims.(*Claims)
	return userClaims, ok
}

// GetUidFromGinContext 从Gin上下文中获取用户ID
func GetUidFromGinContext(c *gin.Context) (int, bool) {
	claims, ok := GetUserFromGinContext(c)
	if !ok {
		return 0, false
	}
	return claims.Uid, true
}

// GetUsernameFromGinContext 从Gin上下文中获取用户名
func GetUsernameFromGinContext(c *gin.Context) (string, bool) {
	claims, ok := GetUserFromGinContext(c)
	if !ok {
		return "", false
	}
	return claims.Username, true
}

// GetUserFromContext 从标准上下文中获取用户声明
func GetUserFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(UserClaimsContextKey).(*Claims)
	return claims, ok
}

// GetUidFromContext 从标准上下文中获取用户ID
func GetUidFromContext(ctx context.Context) (int, bool) {
	claims, ok := GetUserFromContext(ctx)
	if !ok {
		return 0, false
	}
	return claims.Uid, true
}

// GetUsernameFromContext 从标准上下文中获取用户名
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	claims, ok := GetUserFromContext(ctx)
	if !ok {
		return "", false
	}
	return claims.Username, true
}
