package middleware

import (
	"net/http"
	"server/internal/product/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleWare(secret []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization"})
			c.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &service.Claims{}, func(token *jwt.Token) (any, error) {
			return secret, nil
		}) //解析token
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization"})
			c.Abort()
			return
		}
		claim, ok := token.Claims.(*service.Claims)
		if ok {
			c.Set("userID", claim.UserID)
			c.Set("account", claim.Account)
		}
		c.Next()
	}
}
