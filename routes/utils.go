package routes

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func responseHandler(c echo.Context, status int, data any, err error, format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	if err != nil {
		log.Error().Msgf("msg: %s, error: %v", msg, err)
		return c.JSON(status, echo.Map{
			"message": msg,
			"error":   err.Error(),
		})
	}

	log.Info().Msg(msg)
	return c.JSON(status, echo.Map{
		"message": msg,
		"data":    data,
	})
}
