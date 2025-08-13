package midd

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/teoit/gosctx/component/jwtc"
	"github.com/teoit/gosctx/core"
)

func AuthMidd(jwtProvider jwtc.JWTProvider, whitelist *map[string]bool, tokenName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pathUrl := c.Path()
		token := c.Cookies(tokenName)
		if whitelist != nil && (*whitelist)[pathUrl] {
			if token != "" {
				referer := c.Get("Referer")
				if referer == "" || !strings.Contains(referer, c.Hostname()) {
					return c.Redirect("/")
				} else {
					return c.Redirect(referer)
				}
			}
			return c.Next()
		}
		if token == "" {
			return c.Redirect("/login")
		}
		data, err := jwtProvider.ParseToken(c.Context(), token)
		if err != nil {
			return c.Redirect("/login")
		}
		c.Locals(core.KeyRequester, core.NewRequester(data.Subject, data.ID))

		return c.Next()
	}
}

/**
 * Middleware to check permission
 */
