package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

func AuthenticateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		header := c.GetHeader("Authorization")

		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "認証ヘッダーが必要です"})
			return
		}

		tokenString := header[len(BearerSchema):]
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("予期しないサインイン方法: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "トークンの解析に失敗しました"})
			return
		}

		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
			c.Set("email", claims.Subject)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無効なトークンです"})
		}
	}
}

func GenerateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

	return tokenString, err
}
