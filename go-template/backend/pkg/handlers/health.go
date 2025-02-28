package handlers

import "github.com/gofiber/fiber/v2"

// Health TODO: normalde db filan kontrol edilmeli...
func Health(c *fiber.Ctx) error {
	type Status struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	//var home model.Home
	//err := database.DB().First(&home).Error
	//if err != nil {
	//	return c.Status(400).JSON(Status{
	//		Status:  "error",
	//		Message: "database objesi alınamadı",
	//	})
	//}

	//err = cachesystem.Ping()
	//if err != nil {
	//	c.JSON(400, Status{
	//		Status:  "error",
	//		Message: "redis ping error",
	//	})
	//	return
	//}

	return c.JSON(Status{
		Status:  "OK",
		Message: "maşşşallah len",
	})
}
