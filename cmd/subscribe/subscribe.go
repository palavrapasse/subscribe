package main

import (
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	as "github.com/palavrapasse/aspirador/pkg"
	"github.com/palavrapasse/subscribe/internal/data"
	"github.com/palavrapasse/subscribe/internal/http"
	"github.com/palavrapasse/subscribe/internal/logging"
)

func main() {

	logging.Aspirador = as.WithClients(logging.CreateAspiradorClients())

	logging.Aspirador.Trace("Starting Subscribe Service")

	e := echo.New()

	defer e.Close()

	leaksdbctx, oerr := data.OpenLeaksDB()

	if oerr != nil {

		logging.Aspirador.Error(fmt.Sprintf("Could not open Leaks DB connection on server start %v", oerr.Error()))

		return
	}

	subscriptionsdbctx, oerr := data.OpenSubscriptionsDB()

	if oerr != nil {

		logging.Aspirador.Error(fmt.Sprintf("Could not open Subscriptions DB connection on server start %v", oerr.Error()))

		return
	}

	http.RegisterMiddlewares(e, leaksdbctx, subscriptionsdbctx)
	http.RegisterHandlers(e)

	e.Logger.Fatal(http.Start(e))
}
