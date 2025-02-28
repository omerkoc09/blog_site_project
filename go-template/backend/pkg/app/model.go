package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	RequestIDKey = "requestid"
	UserJwtKey   = "user"    // jwt middlewarei ctx.Localse koyuyor
	UserIDKey    = "uid"     // jwt claimsde kullanıyoruz
	UserUUIDKey  = "user_id" // jwt claimsde kullanıyoruz
)

func RequestID(c *fiber.Ctx) string {
	id := c.Locals(RequestIDKey)
	if str, ok := id.(string); ok {
		return str
	}
	return ""
}

func GetUserID(c *fiber.Ctx) int64 {
	claims := GetClaims(c)
	if claims == nil {
		return 0
	}
	userID, ok := claims[UserIDKey].(float64)
	if !ok {
		return 0
	}
	return int64(userID)
}

func GetUserUUID(c *fiber.Ctx) uuid.UUID {
	claims := GetClaims(c)
	if claims == nil {
		return uuid.Nil
	}
	strID, ok := claims[UserUUIDKey].(string)
	if !ok {
		return uuid.Nil
	}

	userID, err := uuid.Parse(strID)
	if err != nil {
		return uuid.Nil
	}
	return userID
}

func GetClaims(c *fiber.Ctx) jwt.MapClaims {
	u := c.Locals(UserJwtKey)
	if u == nil {
		return nil
	}
	user, ok := u.(*jwt.Token)
	if !ok {
		return nil
	}
	return user.Claims.(jwt.MapClaims)
}
