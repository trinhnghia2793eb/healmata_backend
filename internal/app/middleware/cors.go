package middleware

import "github.com/gin-gonic/gin"

func CORS(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// get origin from request header
		origin := c.Request.Header.Get("Origin")
		allowOrigin := ""

		// if config is "*" ---> allow all origin (dev only)
		if len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
			allowOrigin = "*"
		} else {
			// check if request origin is in whitelist origin
			for _, o := range allowedOrigins {
				if o == origin {
					allowOrigin = origin
					break
				}
			}
		}
		// if allowOrigin have value --> valid
		if allowOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		// if it is Request Pre-flight (OPTIONS)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
