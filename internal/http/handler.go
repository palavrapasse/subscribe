package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/palavrapasse/subscribe/internal/data"
	"github.com/palavrapasse/subscribe/internal/logging"
)

func RegisterHandlers(e *echo.Echo) {

	e.POST(subscribeRoute, SubscribeToLeaks)
	e.POST(notificationRoute, NotificationOfLeaks)

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

	err := data.StoreSubscriptionDB(mwctx.SubscriptionsDB, subscription)

	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while storing subscription into DB: %s", err))

		return InternalServerError(ectx)
	}

	logging.Aspirador.Trace("Success in subscribing to leaks.")

	return NoContent(ectx)
}

func NotificationOfLeaks(ectx echo.Context) error {

	logging.Aspirador.Trace("Notification of new leaks")

	request := NotificationRequest{}
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

	querySubscriptionResult, err := data.QuerySubscriptionsDB(mwctx.SubscriptionsDB)

	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while quering subscription from DB: %s", err))

		return InternalServerError(ectx)
	}

	// TODO: delete this once we integrate email send
	logMessage := "\n"
	for _, v := range querySubscriptionResult {
		logMessage += string(v.Subscriber.B64Email) + "\n"

		for _, i := range v.Affected {
			logMessage += "\t\t" + string(i.HSHA256Email) + "\n"
		}
	}
	logging.Aspirador.Trace(logMessage)

	logging.Aspirador.Trace("Success in notification of new leaks.")

	return NoContent(ectx)
}

func useNotFoundHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	}
}
