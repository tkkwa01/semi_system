package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"semi_systems/config"
	"strings"
)

func Auth(must bool, session bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string
		if session {
			// セッションからトークンを取得
			sess := sessions.Default(c)
			token := sess.Get("user")
			if t, ok := token.(string); ok {
				c.Request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t))
			}
		}

		if tokenString == "" {
			tokenString = extractToken(c)
		}

		if tokenString == "" {
			// トークンが見つからない場合
			if must {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.Next()
			}
			return
		}

		// トークンの検証
		claims, err := verifyToken(tokenString, config.Env.App.Secret)
		if err != nil {
			if must {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.Next()
			}
			return
		}

		if userID, ok := claims["user_id"].(float64); ok {
			c.Set("user_id", userID)
		}

		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	bearerToken := c.GetHeader("Authorization")
	if len(bearerToken) > 7 && strings.ToUpper(bearerToken[:7]) == "BEARER " {
		return bearerToken[7:]
	}
	return ""
}

func verifyToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
