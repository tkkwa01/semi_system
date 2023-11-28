package middleware

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"semi_systems/config"
	"strconv"
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
		log.Printf("Verifying token: %s\n", tokenString)
		claims, err := verifyToken(tokenString, config.Env.App.Secret)
		if err != nil {
			if must {
				c.AbortWithStatus(http.StatusUnauthorized)
			} else {
				c.Next()
			}
			return
		}

		if userIDStr, ok := claims["uid"].(string); ok {
			userID, err := strconv.Atoi(userIDStr)
			if err != nil {
				log.Printf("Error converting userID to int: %v", err)
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			c.Set("user_id", uint(userID))
		}

		if userName, ok := claims["user_name"].(string); ok {
			c.Set("user_name", userName)
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
