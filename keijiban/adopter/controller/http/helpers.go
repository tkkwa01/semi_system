package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func bind(c *gin.Context, request interface{}) (ok bool) {
	if err := c.BindJSON(request); err != nil {
		c.Status(http.StatusBadRequest)
		return false
	} else {
		return true
	}
}

// stringをuintに変換する
// Queryにしたらstringになっちゃうからです
func stringToUint(s string) (uint, error) {
	var n uint
	for _, c := range []byte(s) {
		if c < '0' || c > '9' {
			return 0, http.ErrNotSupported
		}
		n = n*10 + uint(c-'0')
	}
	return n, nil
}
