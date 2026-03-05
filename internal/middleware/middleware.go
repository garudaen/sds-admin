package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS handles CORS headers for cross-origin requests.
// CORS 处理跨域请求的CORS头部
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// Logger logs request information.
// Logger 记录请求信息
func Logger() gin.HandlerFunc {
	return gin.Logger()
}

// Recovery recovers from panics and returns a 500 error.
// Recovery 从panic中恢复并返回500错误
func Recovery() gin.HandlerFunc {
	return gin.Recovery()
}
