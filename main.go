package main

import (
	"github.com/italosilva18/funeraria/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Rotas para Admin
	e.POST("/admin", handlers.CreateAdmin)

	// Aqui você adicionaria mais rotas para funcionario, cliente, produto e plano

	e.Logger.Fatal(e.Start(":1323"))
}
