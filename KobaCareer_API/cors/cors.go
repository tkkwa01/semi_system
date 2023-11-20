package cors

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CorsOption struct {
	AllowOrigins []string
	MaxAge       time.Duration
}

func Cors(option *CorsOption) echo.MiddlewareFunc {
	var (
		allowOrigins []string
		maxAge       = time.Hour * 2
	)

	if option != nil {
		allowOrigins = option.AllowOrigins

		if option.MaxAge != 0 {
			maxAge = option.MaxAge
		}
	}

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Length", "Content-Type", "X-Request-Id"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead},
		MaxAge:           int(maxAge.Seconds()),
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) (bool, error) {
			if len(allowOrigins) == 0 {
				return true, nil
			}
			for _, allowOrigin := range allowOrigins {
				if allowOrigin == origin {
					return true, nil
				}
			}
			return false, nil
		},
	})
}
