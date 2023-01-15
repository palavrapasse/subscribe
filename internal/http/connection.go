package http

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

const (
	serverHostEnvKey = "server_host"
	serverPortEnvKey = "server_port"
)

var (
	serverHost = os.Getenv(serverHostEnvKey)
	serverPort = os.Getenv(serverPortEnvKey)
)

func Start(e *echo.Echo) error {
	return e.Start(serverAddress())
}

func serverAddress() string {
	return fmt.Sprintf("%s:%s", serverHost, serverPort)
}
