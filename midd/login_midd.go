package midd

import (
	"github.com/teoit/gosctx/configs"
	"github.com/teoit/gosctx/core"

	"github.com/gofiber/fiber/v2"
)

/**
 * Middleware check login
 */
func LoginMidd(jwtProvider configs.JWTProvider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("jwtToken")
		if token == "" {
			return c.Redirect("/login")
		}
		data, err := jwtProvider.ParseToken(c.Context(), token)
		if err != nil {
			return c.Redirect("/login")
		}

		c.Locals("requester", core.NewRequester(data.Subject, data.ID))

		return c.Next()
	}
}
