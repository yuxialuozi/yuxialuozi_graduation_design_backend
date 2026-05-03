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

		if authHeader == "" {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(parts[1], cfg.JWT.Secret)
		if err != nil {
			response.Unauthorized(c, "无效的 token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("tenantId", claims.TenantID)

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

func GetTenantID(c *gin.Context) uint {
	tenantID, exists := c.Get("tenantId")
	if !exists {
		return 0
	}
	return tenantID.(uint)
}

// RequireRole 创建角色验证中间件，要求用户具有指定角色
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := GetRole(c)
		if userRole == "" {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}
		if userRole != role {
			response.Forbidden(c, "没有权限访问此资源")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireAdmin 要求用户具有admin角色
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// RequireUser 要求用户具有user角色（租户端用户）
func RequireUser() gin.HandlerFunc {
	return RequireRole("user")
}
