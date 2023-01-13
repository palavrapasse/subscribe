package http

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/palavrapasse/damn/pkg/database"
	"github.com/palavrapasse/subscribe/internal/logging"
)

const (
	dbMiddlewareKey = "db"
)

type MiddlewareContext struct {
	DB database.DatabaseContext[database.Record]
}

func RegisterMiddlewares(e *echo.Echo, dbctx database.DatabaseContext[database.Record]) {
	e.Use(dbAccessMiddleware(dbctx))
	e.Use(loggingMiddleware())
}

func GetMiddlewareContext(ectx echo.Context) (MiddlewareContext, error) {
	db, dok := ectx.Get(dbMiddlewareKey).(database.DatabaseContext[database.Record])
	var err error

	if !dok {
		err = errors.New("DB not available in middleware")
	}

	return MiddlewareContext{
		DB: db,
	}, err
}

func dbAccessMiddleware(dbctx database.DatabaseContext[database.Record]) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ectx echo.Context) error {
			ectx.Set(dbMiddlewareKey, dbctx)

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
