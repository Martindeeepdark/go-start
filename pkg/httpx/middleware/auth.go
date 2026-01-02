package middleware

import (
	"strings"

	"github.com/Martindeeepdark/go-start/pkg/commonadapter"
	"github.com/Martindeeepdark/go-start/pkg/httpx/response"
	"github.com/gin-gonic/gin"
)

const userIDKey = "UserID"

// RequireAuth 认证中间件
// 从 Authorization 头提取 Bearer Token，调用能力位完成校验，并在上下文设置 UserID。
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) > 7 && (strings.HasPrefix(token, "Bearer ") || strings.HasPrefix(token, "bearer ")) {
			token = token[7:]
		}
		auth, _, _, _, _, _ := commonadapter.Abilities()
		userID, err := auth.VerifyToken(token)
		if err != nil || userID == "" {
			response.Error(c, 401, "未授权")
			c.Abort()
			return
		}
		c.Set(userIDKey, userID)
		c.Next()
	}
}

// RequirePermission 权限校验中间件
// 依赖 RequireAuth 已设置的 UserID，或备用令牌校验，随后校验指定权限码。
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userID string
		if v, ok := c.Get(userIDKey); ok {
			if s, ok2 := v.(string); ok2 {
				userID = s
			}
		}
		if userID == "" {
			token := c.GetHeader("Authorization")
			if len(token) > 7 && (strings.HasPrefix(token, "Bearer ") || strings.HasPrefix(token, "bearer ")) {
				token = token[7:]
			}
			auth, _, _, _, _, _ := commonadapter.Abilities()
			uid, err := auth.VerifyToken(token)
			if err != nil || uid == "" {
				response.Error(c, 401, "未授权")
				c.Abort()
				return
			}
			userID = uid
		}
		auth, _, _, _, _, _ := commonadapter.Abilities()
		if err := auth.RequirePermission(userID, permission); err != nil {
			response.Error(c, 403, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}
