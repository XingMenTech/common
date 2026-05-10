package utils

// JWT认证组件使用示例

// 示例1: 生成令牌对
func ExampleGenerateTokenPair() {
	// 生成访问令牌和刷新令牌
	accessToken, refreshToken, err := GenerateTokenPair(
		123,               // 用户ID
		"username",        // 用户名
		1,                 // 用户类型
		"device_id_12345", // 设备ID
		"192.168.1.100",   // 客户端IP
	)
	if err != nil {
		// 处理错误
		panic(err)
	}

	// accessToken 和 refreshToken 可以返回给客户端
	println("Access Token:", accessToken)
	println("Refresh Token:", refreshToken)
}

// 示例2: 验证访问令牌
func ExampleValidateAccessToken() {
	token := "your_jwt_token_here"

	claims, err := ValidateAccessToken(token)
	if err != nil {
		// 令牌无效或已过期
		println("Invalid token:", err.Error())
		return
	}

	// 访问令牌有效，可以使用claims中的信息
	println("User ID:", claims.Uid)
	println("Username:", claims.Username)
	println("User Type:", claims.UserType)
}

// 示例3: 使用刷新令牌获取新的访问令牌
func ExampleRefreshAccessToken() {
	refreshToken := "your_refresh_token_here"

	newAccessToken, newRefreshToken, err := RefreshAccessToken(refreshToken)
	if err != nil {
		// 刷新令牌无效或已过期
		println("Failed to refresh token:", err.Error())
		return
	}

	// 使用新的令牌对
	println("New Access Token:", newAccessToken)
	println("New Refresh Token:", newRefreshToken)
}

// 示例4: 在Gin框架中使用JWT中间件
func ExampleGinMiddleware() {
	/*
		import (
			"github.com/gin-gonic/gin"
			"github.com/XingMenTech/common/utils"
		)

		router := gin.Default()

		// 创建JWT中间件
		jwtMiddleware := utils.NewJWTMiddleware()

		// 可选：设置排除路径（不需要认证的路径）
		jwtMiddleware.ExcludedPaths = []string{"/login", "/register", "/public"}

		// 可选：设置自定义验证器
		jwtMiddleware.CustomValidator = func(claims *utils.Claims) error {
			// 例如：检查用户是否有特定权限
			if claims.UserType != 1 {
				return errors.New("insufficient permissions")
			}
			return nil
		}

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
					"device_id":  claims.DeviceId,
					"client_ip":  claims.ClientIp,
				})
			})
		}

		router.Run(":8080")
	*/
}

// 示例5: 从令牌中提取用户信息
func ExampleGetUserInfoFromToken() {
	token := "your_jwt_token_here"

	// 获取用户ID
	uid, err := GetUidFromToken(token)
	if err != nil {
		println("Failed to get UID:", err.Error())
		return
	}
	println("User ID:", uid)

	// 获取用户名
	username, err := GetUsernameFromToken(token)
	if err != nil {
		println("Failed to get username:", err.Error())
		return
	}
	println("Username:", username)
}

// 示例6: 单独生成访问令牌和刷新令牌
func ExampleGenerateSeparateTokens() {
	// 生成访问令牌
	accessToken, err := GenerateAccessToken(
		123,               // 用户ID
		"username",        // 用户名
		1,                 // 用户类型
		"device_id_12345", // 设备ID
		"192.168.1.100",   // 客户端IP
	)
	if err != nil {
		panic(err)
	}

	// 解析访问令牌以获取JTI
	claims, err := ParseAccessToken(accessToken)
	if err != nil {
		panic(err)
	}

	// 生成刷新令牌
	refreshToken, err := GenerateRefreshToken(claims.Uid, claims.Id)
	if err != nil {
		panic(err)
	}

	println("Access Token:", accessToken)
	println("Refresh Token:", refreshToken)
}
