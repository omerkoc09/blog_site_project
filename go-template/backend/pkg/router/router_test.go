package router

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"github.com/hayrat/go-template2/backend/pkg/app"
)

func TestCtxWrap(t *testing.T) {
	f := func(ctx *app.Ctx) error {
		return nil
	}
	r := ctxWrap(f)
	fb := fiber.New()
	ctx := fb.AcquireCtx(&fasthttp.RequestCtx{})

	err := r(ctx)
	assert.Nil(t, err)
}

func h(ctx *app.Ctx) error {
	return nil
}

const testRoute = "/test"

func TestGet(t *testing.T) {
	fb := fiber.New()
	Get(fb, testRoute, h)
	req := httptest.NewRequest("GET", testRoute, nil)
	resp, err := fb.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	req = httptest.NewRequest("GET", "/asdf", nil)
	resp, err = fb.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestPost(t *testing.T) {
	fb := fiber.New()
	Post(fb, testRoute, h)
	req := httptest.NewRequest("POST", testRoute, nil)
	resp, err := fb.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	req = httptest.NewRequest("POST", "/asdf", nil)
	resp, err = fb.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestPut(t *testing.T) {
	fb := fiber.New()
	Put(fb, testRoute, h)
	req := httptest.NewRequest("PUT", testRoute, nil)
	resp, err := fb.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	req = httptest.NewRequest("PUT", "/asdf", nil)
	resp, err = fb.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestDelete(t *testing.T) {
	fb := fiber.New()
	Delete(fb, testRoute, h)
	req := httptest.NewRequest("DELETE", testRoute, nil)
	resp, err := fb.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	req = httptest.NewRequest("DELETE", "/asdf", nil)
	resp, err = fb.Test(req)
	assert.Nil(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}
