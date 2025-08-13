package midd

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/teoit/gosctx"
)

type CanGetStatusCode interface {
	StatusCode() int
}

func Recovery(serviceCtx gosctx.ServiceContext) fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				// panic(err)
				// c.Set("Content-Type", "application/json")
				if appErr, ok := err.(CanGetStatusCode); ok {
					if err = c.Status(appErr.StatusCode()).JSON(&fiber.Map{
						"errors": appErr,
					}); err != nil {
						return
					}
				} else {
					if err := c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
						"code":    http.StatusInternalServerError,
						"status":  "internal server error",
						"message": "something went wrong, please try again or contact supporters",
					}); err != nil {
						return
					}
				}
				serviceCtx.Logger("service").Errorf("%+v \n", err)
			}
		}()
		return c.Next()
	}
}
