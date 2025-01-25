package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Token is required",
			})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid token format",
			})
			ctx.Abort()
			return
		}

		claims, err := VerifyToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid or expired token",
				"error":   err.Error(),
			})
			ctx.Abort()
			return
		}

		ctx.Set("role", claims.Role)

		userRoleInterface, exists := ctx.Get("role")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Role not found in context",
			})
			ctx.Abort()
			return
		}

		userRole, ok := userRoleInterface.(string)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Invalid role type in context",
			})
			ctx.Abort()
			return
		}

		allowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			ctx.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": "You do not have permission to access this resource",
				"details": gin.H{
					"required_roles": allowedRoles,
					"your_role":      userRole,
				},
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
