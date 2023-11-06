package api

import (
	"funeraria/api/handlers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	// Rotas para produtos
	e.POST("/products", handlers.CreateProduct)
	e.GET("/products/:id", handlers.GetProduct)
	e.PUT("/products/:id", handlers.UpdateProduct)
	e.DELETE("/products/:id", handlers.DeleteProduct)

	// Rotas para planos (você deve criar manipuladores de plano semelhantes)
	e.POST("/plans", handlers.CreatePlan)
	e.GET("/plans/:id", handlers.GetPlan)
	e.PUT("/plans/:id", handlers.UpdatePlan)
	e.DELETE("/plans/:id", handlers.DeletePlan)
}
