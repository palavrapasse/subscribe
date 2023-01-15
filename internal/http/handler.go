package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/palavrapasse/subscribe/internal/logging"
)

func RegisterHandlers(e *echo.Echo) {

	e.POST(subscribeRoute, SubscribeToLeaks)

	echo.NotFoundHandler = useNotFoundHandler()
}

func SubscribeToLeaks(ectx echo.Context) error {

	logging.Aspirador.Trace("Subscribing to leaks")

	request := SubscriptionRequest{}
	decerr := ectx.Bind(&request)

	if decerr != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while reading request body: %s", decerr))

		return InternalServerError(ectx)
	}

	mwctx, gmerr := GetMiddlewareContext(ectx)

	if gmerr != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while getting Middleware Context: %s", gmerr))

		return InternalServerError(ectx)
	}

	subscription := SubscriptionRequestToSubscription(request)

	err := mwctx.SubscriptionsDB.InsertSubscription(subscription)

	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while storing subscription into DB: %s", err))

		return InternalServerError(ectx)
	}

	logging.Aspirador.Trace("Success in subscribing to leaks.")

	return NoContent(ectx)
}

func useNotFoundHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	}
}
