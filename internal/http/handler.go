package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterHandlers(e *echo.Echo) {

	e.POST(subscribeRoute, SubscribeToLeaks)

	echo.NotFoundHandler = useNotFoundHandler()
}

func SubscribeToLeaks(ectx echo.Context) error {

	return Ok(ectx, nil)
}

func useNotFoundHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	}
}
