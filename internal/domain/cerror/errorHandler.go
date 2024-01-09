package cerror

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {

	if err != nil {
		if errors.Is(err, ErrUserExists) {
			err := fmt.Sprintf("login is already exists")
			return c.Status(409).JSON(fiber.Map{
				"Message": err,
			})
		}
		if errors.Is(err, ErrAppExists) {
			err := fmt.Sprintf("app is already exists")
			return c.Status(409).JSON(fiber.Map{
				"Message": err,
			})
		}
		if errors.Is(err, ErrInvalidCredentials) {
			err := fmt.Sprintf("incorrect data")
			return c.Status(400).JSON(fiber.Map{
				"Message": err,
			})
		}
		if errors.Is(err, ErrAppNotFound) {
			err := fmt.Sprintf("app not found")
			return c.Status(400).JSON(fiber.Map{
				"Message": err,
			})
		}

		err := fmt.Sprintf("internal error server")
		return c.Status(500).JSON(fiber.Map{
			"Message": err,
		})
	}
	return c.Status(200).SendString("success")
}
