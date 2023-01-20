package http

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/palavrapasse/damn/pkg/entity"
	"github.com/palavrapasse/damn/pkg/entity/query"
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

	leakId := entity.AutoGenKey(request.LeakId)
	usersAffectedByLeak, err := getUsersAffectedByLeak(mwctx, leakId)

	if err != nil {
		return InternalServerError(ectx)
	}

	// TODO: delete this once we integrate email send
	var logMessage = "\n"
	for _, v := range usersAffectedByLeak {
		email, _ := base64.StdEncoding.DecodeString(string(v.DestinationB64Email))
		logMessage += "Email to: " + string(v.DestinationB64Email) + " -> " + string(email) + "\n"
		logMessage += "\t\t The affecteds by leak are:\n"
		for _, i := range v.AffectedsEmail {
			logMessage += "\t\t\t- " + string(i) + "\n"
		}
	}
	logging.Aspirador.Trace(logMessage)

	leakInformation, err := getLeakInformation(mwctx, leakId)

	if err != nil {
		return InternalServerError(ectx)
	}

	emailInfo := data.EmailInfo{
		UsersAffected: usersAffectedByLeak,
		LeakInfo:      leakInformation,
	}

	logMessage = "\nLeak.Context \t" + string(emailInfo.LeakInfo.Leak.Context)
	logMessage += "\nLeak.ShareDateSC " + emailInfo.LeakInfo.Leak.ShareDateSC.String()
	for _, v := range emailInfo.LeakInfo.PlatformsAffected {
		logMessage += "\n\t\tPlatform " + v.Name
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

func getUsersAffectedByLeak(mwctx MiddlewareContext, leakId entity.AutoGenKey) (data.AllAffectedsInfo, error) {
	querySubscriptionResult, err := data.QuerySubscriptionsDB(mwctx.SubscriptionsDB)

	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while quering subscription from DB: %s", err))

		return nil, err
	}

	affectedSubscription := querySubscriptionResult.GetAffectUsers()

	usersAffectedByLeak, err := data.QueryAffectedLeaksDB(mwctx.LeaksDB, leakId, affectedSubscription)

	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while quering leak affected users from DB: %s", err))

		return nil, err
	}

	subscriberAffectedsByLeak := querySubscriptionResult.GetAffectedsInfo(usersAffectedByLeak)

	return subscriberAffectedsByLeak, nil
}

func getLeakInformation(mwctx MiddlewareContext, leakId entity.AutoGenKey) (data.LeakInfo, error) {
	queryLeakResult, err := data.QueryLeakByIdDB(mwctx.LeaksDB, leakId)

	if err != nil {
		logging.Aspirador.Error(fmt.Sprintf("Error while quering leak from DB: %s", err))

		return data.LeakInfo{}, err
	}

	size := len(queryLeakResult)
	if size < 1 {
		err = fmt.Errorf("Query to get leak information did not return a thing")
		logging.Aspirador.Error(fmt.Sprintf("Error while quering leak from DB: %s", err))
		return data.LeakInfo{}, err
	}

	leak := query.Leak{
		Context:     queryLeakResult[0].Context,
		ShareDateSC: queryLeakResult[0].ShareDateSC,
		LeakId:      queryLeakResult[0].LeakId,
	}

	platforms := make([]query.Platform, size)
	for i, v := range queryLeakResult {
		platforms[i] = v.Platform
	}

	return data.LeakInfo{
		Leak:              leak,
		PlatformsAffected: platforms,
	}, nil
}
