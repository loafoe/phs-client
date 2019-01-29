package client

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GinAuthFilter(sharedSecret string) gin.HandlerFunc {
	verifier := NewSignatureVerifier(sharedSecret)

	return func(c *gin.Context) {
		if valid, err := verifier.VerifyRequest(c.Request); !valid {
			c.AbortWithError(http.StatusUnauthorized, err)
		}
		c.Next()
	}
}
