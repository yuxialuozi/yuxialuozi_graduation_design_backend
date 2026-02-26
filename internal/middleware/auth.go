package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"yuxialuozi_graduation_design_backend/internal/config"
	"yuxialuozi_graduation_design_backend/pkg/response"
	"yuxialuozi_graduation_design_backend/pkg/utils"
)

func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		// 添加调试日志
		method := c.Request.Method
		path := c.Request.URL.Path

		if authHeader == "" {
			// 记录详细错误日志
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// 记录授权头信息（不记录完整token）
		authSnippet := ""
		if len(authHeader) > 20 {
			authSnippet = authHeader[:20] + "..."
		} else {
			authSnippet = authHeader
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1], cfg.JWT.Secret)
		if err != nil {
			// 记具体错误
			response.Unauthorized(c, "无效的 token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func GetUserID(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	return userID.(uint)
}

func GetUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}

func GetRole(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}
	return role.(string)
}
