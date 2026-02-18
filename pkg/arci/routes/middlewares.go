package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{
				"error": "authorization header required",
			})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(401, gin.H{
				"error": "invalid token claims",
			})
			c.Abort()
			return
		}

		c.Set("member_id", claims.MemberID)
		c.Set("showname", claims.ShowName)
		c.Set("is_admin", claims.IsAdmin)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			c.JSON(403, gin.H{
				"error": "admin access required",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
