package http

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/palavrapasse/damn/pkg/database"
	"github.com/palavrapasse/subscribe/internal/logging"
)

const (
	leaksDBMiddlewareKey         = "db_leak"
	subscriptionsDBMiddlewareKey = "db_subscriptions"
)

type MiddlewareContext struct {
	LeaksDB         database.DatabaseContext[database.Record]
	SubscriptionsDB database.DatabaseContext[database.Record]
}

func RegisterMiddlewares(e *echo.Echo, leakdbctx database.DatabaseContext[database.Record], subscritpiondbctx database.DatabaseContext[database.Record]) {
	e.Use(dbAccessMiddleware(leakdbctx, leaksDBMiddlewareKey))
	e.Use(dbAccessMiddleware(subscritpiondbctx, subscriptionsDBMiddlewareKey))
	e.Use(loggingMiddleware())
}

func GetMiddlewareContext(ectx echo.Context) (MiddlewareContext, error) {
	db, dok := ectx.Get(leaksDBMiddlewareKey).(database.DatabaseContext[database.Record])
	var err error

	if !dok {
		err = errors.New("leaks DB not available in middleware")
		return MiddlewareContext{}, err
	}

	dbSubscribe, dok := ectx.Get(subscriptionsDBMiddlewareKey).(database.DatabaseContext[database.Record])

	if !dok {
		err = errors.New("subscriptions DB not available in middleware")
	}

	return MiddlewareContext{
		LeaksDB:         db,
		SubscriptionsDB: dbSubscribe,
	}, err
}

func dbAccessMiddleware(dbctx database.DatabaseContext[database.Record], middlewareKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			ectx.Set(middlewareKey, dbctx)

			return next(ectx)
		}
	}
}

func loggingMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {

			req := ectx.Request()

			logging.Aspirador.Info(fmt.Sprintf("Host: %s | Method: %s | Path: %s", req.Host, req.Method, req.URL.RequestURI()))

			return next(ectx)
		}
	}
}
