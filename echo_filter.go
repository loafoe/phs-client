package client

import (
	"github.com/labstack/echo"
)

func EchoAuthFilter(sharedSecret string) echo.MiddlewareFunc {
	verifier := NewSignatureVerifier(sharedSecret)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if valid, _ := verifier.VerifyRequest(c.Request()); !valid {
				return echo.ErrUnauthorized
			}
			return next(c)
		}
	}
}
