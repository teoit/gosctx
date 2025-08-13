package core

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func WriteErrorResponse(c *fiber.Ctx, err error) error {
	if errSt, ok := err.(StatusCodeCarrier); ok {
		return c.Status(errSt.StatusCode()).JSON(errSt)
	}
	return c.Status(http.StatusInternalServerError).JSON(ErrInternalServerError.WithError(err.Error()))
}

func ReturnErrsForApi(ctx *fiber.Ctx, msgErr interface{}) error {
	return WriteErrorResponse(ctx, ErrBadRequest.WithDetail("msg", msgErr))
}

func ReturnErrForApi(ctx *fiber.Ctx, msgErr string) error {
	return WriteErrorResponse(ctx, ErrBadRequest.WithDetail("msg", msgErr))
}
