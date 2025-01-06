package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// RequireAuth checks for a valid JWT token
func RequireAuth(c *fiber.Ctx) error {
	tokenString := c.Cookies("jwt_token") // Try to get token from cookies
	if tokenString == "" {
		tokenString = c.Get("Authorization") // Or from Authorization header
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized: No token provided")
		}
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// ensure the signing method is valid
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized: Invalid token")
	}

	// Token is valid, allow the request to proceed
	return c.Next()
}
