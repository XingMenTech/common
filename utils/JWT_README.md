# JWT 认证组件

## 概述

这是一个完整的JWT（JSON Web Token）认证组件，提供了安全的身份验证和授权功能。该组件支持访问令牌和刷新令牌机制，适用于Web应用和API服务。

## 特性

- ✅ 标准的JWT实现，基于 `github.com/golang-jwt/jwt/v4`
- ✅ 访问令牌和刷新令牌双令牌机制
- ✅ 可配置的令牌过期时间
- ✅ 安全的签名验证
- ✅ 丰富的用户信息存储（UID、用户名、用户类型、设备ID、客户端IP）
- ✅ Gin框架中间件支持
- ✅ 防止重放攻击（JTI）
- ✅ 完整的单元测试覆盖

## 安装

该组件已包含在 `github.com/XingMenTech/common` 模块中，无需额外安装。

确保你的项目中包含以下依赖：

```go
github.com/golang-jwt/jwt/v4 v4.5.2
github.com/gin-gonic/gin v1.10.0  // 如果使用Gin中间件
```

## 快速开始

### 1. 生成令牌对

```go
import "github.com/XingMenTech/common/utils"

// 生成访问令牌和刷新令牌
accessToken, refreshToken, err := utils.GenerateTokenPair(
    123,                    // 用户ID
    "username",             // 用户名
    1,                      // 用户类型
    "device_id_12345",      // 设备ID
    "192.168.1.100",        // 客户端IP
)
if err != nil {
    // 处理错误
    log.Fatal(err)
}

// 将令牌返回给客户端
// accessToken: 用于API请求认证（有效期24小时）
// refreshToken: 用于刷新access token（有效期7天）
```

### 2. 验证访问令牌

```go
// 从请求头中获取token
token := "Bearer your_jwt_token_here"
token = strings.TrimPrefix(token, "Bearer ")

// 验证令牌
claims, err := utils.ValidateAccessToken(token)
if err != nil {
    // 令牌无效或已过期
    log.Println("Invalid token:", err.Error())
    return
}

// 使用claims中的用户信息
log.Println("User ID:", claims.Uid)
log.Println("Username:", claims.Username)
```

### 3. 刷新令牌

```go
// 当access token过期时，使用refresh token获取新的令牌对
newAccessToken, newRefreshToken, err := utils.RefreshAccessToken(refreshToken)
if err != nil {
    // 刷新失败，需要重新登录
    log.Println("Failed to refresh token:", err.Error())
    return
}

// 使用新的令牌对
```

## 在Gin框架中使用

### 基本用法

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/XingMenTech/common/utils"
)

func main() {
    router := gin.Default()

    // 创建JWT中间件
    jwtMiddleware := utils.NewJWTMiddleware()
    
    // 设置排除路径（不需要认证的路径）
    jwtMiddleware.ExcludedPaths = []string{"/login", "/register", "/public"}

    // 应用中间件到路由组
    protected := router.Group("/api/protected")
    protected.Use(jwtMiddleware.GinMiddleware())
    {
        protected.GET("/profile", func(c *gin.Context) {
            // 从上下文中获取用户信息
            claims, exists := utils.GetUserFromGinContext(c)
            if !exists {
                c.JSON(401, gin.H{"error": "unauthorized"})
                return
            }

            c.JSON(200, gin.H{
                "user_id":    claims.Uid,
                "username":   claims.Username,
                "user_type":  claims.UserType,
            })
        })
    }

    router.Run(":8080")
}
```

### 自定义验证器

```go
// 设置自定义验证器
jwtMiddleware.CustomValidator = func(claims *utils.Claims) error {
    // 例如：检查用户是否有特定权限
    if claims.UserType != 1 {
        return errors.New("insufficient permissions")
    }
    
    // 例如：检查设备ID是否匹配
    if claims.DeviceId != expectedDeviceId {
        return errors.New("device not recognized")
    }
    
    return nil
}
```

### 便捷函数

```go
// 从Gin上下文中快速获取用户信息
uid, exists := utils.GetUidFromGinContext(c)
username, exists := utils.GetUsernameFromGinContext(c)
```

## API 参考

### 核心函数

#### GenerateAccessToken
生成访问令牌
```go
func GenerateAccessToken(uid int, username string, userType int, deviceId string, clientIp string) (string, error)
```

#### GenerateRefreshToken
生成刷新令牌
```go
func GenerateRefreshToken(uid int, originalJti string) (string, error)
```

#### GenerateTokenPair
生成访问令牌和刷新令牌对
```go
func GenerateTokenPair(uid int, username string, userType int, deviceId string, clientIp string) (accessToken string, refreshToken string, err error)
```

#### ValidateAccessToken
验证访问令牌
```go
func ValidateAccessToken(token string) (*Claims, error)
```

#### ValidateRefreshToken
验证刷新令牌
```go
func ValidateRefreshToken(token string) (*RefreshClaims, error)
```

#### RefreshAccessToken
使用刷新令牌获取新的访问令牌
```go
func RefreshAccessToken(refreshTokenStr string) (newAccessToken string, newRefreshToken string, err error)
```

### Claims 结构

```go
type Claims struct {
    Uid        int    `json:"uid"`           // 用户ID
    Username   string `json:"username"`      // 用户名
    UserType   int    `json:"user_type"`     // 用户类型
    DeviceId   string `json:"device_id"`     // 设备ID
    ClientIp   string `json:"client_ip"`     // 客户端IP
    jwt.StandardClaims                      // 标准JWT声明
}
```

## 配置

可以在 `Jwt.go` 文件中修改以下常量：

```go
const (
    TokenSecret      = "game" // 加密密钥，建议在生产环境中使用更复杂的密钥
    TokenInvalidTime = 24     // 访问令牌有效期（小时）
    RefreshTokenTime = 7 * 24 // 刷新令牌有效期（小时）
)
```

**安全建议：**
- 在生产环境中，使用环境变量或配置文件管理 `TokenSecret`
- 使用至少32个字符的随机字符串作为密钥
- 定期轮换密钥

## 安全最佳实践

1. **HTTPS传输**：始终通过HTTPS传输令牌
2. **安全存储**：客户端应安全存储令牌（如HttpOnly Cookie）
3. **短期访问令牌**：访问令牌有效期不宜过长（默认24小时）
4. **刷新令牌轮换**：每次刷新时生成新的刷新令牌
5. **令牌撤销**：实现令牌黑名单机制以支持登出功能
6. **速率限制**：对刷新令牌接口实施速率限制
7. **设备绑定**：将令牌与设备ID绑定，防止令牌被盗用

## 测试

运行单元测试：

```bash
cd /Users/zhangyuan/workspase/xm/common
go test ./utils -v -run ".*Token.*"
```

## 示例代码

查看 `jwt_example.go` 文件获取更多使用示例。

## 常见问题

### Q: 如何实现令牌注销（黑名单）？
A: 可以将令牌的JTI存储在Redis中，设置过期时间与令牌剩余有效期一致。在验证令牌时，检查JTI是否在黑名单中。

### Q: 如何支持多设备登录？
A: 当前实现已支持，每个设备会生成独立的令牌对。可以通过DeviceId区分不同设备。

### Q: 如何强制用户下线？
A: 将该用户的所有活跃令牌JTI加入黑名单，或者更改TokenSecret使所有现有令牌失效。

## 许可证

本项目遵循项目主许可证。

## 贡献

欢迎提交Issue和Pull Request来改进这个组件。
