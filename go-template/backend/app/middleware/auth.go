package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hayrat/go-template2/backend/common/service"
)

var tokenParser service.TokenParser

// TokenParser'ı yapılandırın (NewTokenService arayüz ile entegre)
func SetTokenParser(parser service.TokenParser) {
	tokenParser = parser
}

func OptionalJWTAuth(c *fiber.Ctx) error {
	token := c.Cookies("jwt")
	if token == "" {
		return c.Next() // Eğer token yoksa devam et
	}

	// JWT doğrulama işlemi
	claims, err := tokenParser.ParseToken(token)
	if err != nil {
		return c.Next() // Token geçersizse devam et
	}

	// Kullanıcı bilgilerini context'e ekle
	c.Locals("userID", claims.UserID)
	return c.Next()
}
