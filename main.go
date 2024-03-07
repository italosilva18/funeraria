package main

import (
	"net/http"

	"github.com/labstack/echo"
)

type HelloWorld struct {
	Message string `json:"message"`
}

func main() {
	e := echo.New()

	// Rota para "Olá, mundo!"
	e.GET("/hello", Greetings)

	// Rota com parâmetro na URL
	e.GET("/hello/:name", GreetingsWithParams)

	// Rota com consulta (query) na URL
	e.GET("/hello-queries", GreetingsWithQuery)

	e.Logger.Fatal(e.Start(":3000"))
}

func Greetings(c echo.Context) error {
	return c.JSON(http.StatusOK, HelloWorld{
		Message: "Olá, mundo!",
	})
}

func GreetingsWithParams(c echo.Context) error {
	params := c.Param("name")
	return c.JSON(http.StatusOK, HelloWorld{
		Message: "Olá, meu nome é " + params,
	})
}

func GreetingsWithQuery(c echo.Context) error {
	query := c.QueryParam("name")
	return c.JSON(http.StatusOK, HelloWorld{
		Message: "Olá, estou usando queries e meu nome é " + query,
	})
}
