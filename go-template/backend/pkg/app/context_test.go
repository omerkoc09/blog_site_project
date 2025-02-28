package app

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

const testRoute = "/test"

func TestCtx_SuccessResponse(t *testing.T) {
	fb := fiber.New()
	fb.Get(testRoute, func(ctx *fiber.Ctx) error {
		c := Ctx{
			Ctx: ctx,
		}
		err := c.SuccessResponse("test")
		assert.Nil(t, err)

		err = c.SuccessResponse(nil, 1)
		assert.Nil(t, err)
		return nil
	})
	req := httptest.NewRequest("GET", testRoute, nil)
	_, err := fb.Test(req)
	assert.Nil(t, err)
}

func TestCtx_Log(t *testing.T) {
	c := Ctx{
		//Logger: log.GetLogger(""),
	}
	c.SetLogFields()
	c.ErrorLog("error")
	c.InfoLog("info")
}
