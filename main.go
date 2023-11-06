package main

import (
	"fmt"
	"funeraria/api"
	"funeraria/config"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Carregue as configurações do arquivo config.yaml
	if err := config.LoadConfig("config.yaml"); err != nil {
		log.Fatalf("Erro ao carregar as configurações: %v", err)
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Configure e registre as rotas da API
	api.RegisterRoutes(e)

	// Inicie o servidor
	port := config.AppConfig.ServerPort
	address := fmt.Sprintf(":%d", port)
	e.Logger.Fatal(e.Start(address))
}
