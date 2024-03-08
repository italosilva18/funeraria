package api

import (
	"github.com/italosilva18/funeraria/api/routes"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start() {
	// Inicializa a aplicação Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Configura as rotas
	routes.SetupRoutes(e)

	// Inicia o servidor Echo
	e.Logger.Fatal(e.Start(":8080"))
}
