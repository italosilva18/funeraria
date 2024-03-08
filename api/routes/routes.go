package routes

import (
	"github.com/italosilva18/funeraria/api/controllers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	// Rota de cadastro de clientes
	e.POST("/clientes", controllers.CreateCliente)
	e.GET("/clientes", controllers.GetClientes)
	e.GET("/clientes/:id", controllers.GetClienteByID)
	e.PUT("/clientes/:id", controllers.UpdateCliente)
	e.DELETE("/clientes/:id", controllers.DeleteCliente)

	// Rota de gestão de serviços funerários
	e.POST("/servicos", controllers.CreateServico)
	e.GET("/servicos", controllers.GetServicos)
	e.GET("/servicos/:id", controllers.GetServicoByID)
	e.PUT("/servicos/:id", controllers.UpdateServico)
	e.DELETE("/servicos/:id", controllers.DeleteServico)

	// Rota de controle de estoque e produtos
	e.POST("/produtos", controllers.CreateProduto)
	e.GET("/produtos", controllers.GetProdutos)
	e.GET("/produtos/:id", controllers.GetProdutoByID)
	e.PUT("/produtos/:id", controllers.UpdateProduto)
	e.DELETE("/produtos/:id", controllers.DeleteProduto)

	// Rota de agendamento e calendário
	e.POST("/agendamentos", controllers.CreateAgendamento)
	e.GET("/agendamentos", controllers.GetAgendamentos)
	e.GET("/agendamentos/:id", controllers.GetAgendamentoByID)
	e.PUT("/agendamentos/:id", controllers.UpdateAgendamento)
	e.DELETE("/agendamentos/:id", controllers.DeleteAgendamento)

	// Rota de financeiro e cobranças
	e.POST("/faturas", controllers.CreateFatura)
	e.GET("/faturas", controllers.GetFaturas)
	e.GET("/faturas/:id", controllers.GetFaturaByID)
	e.PUT("/faturas/:id", controllers.UpdateFatura)
	e.DELETE("/faturas/:id", controllers.DeleteFatura)

	// Rota de autenticação do administrador
	e.POST("/admin/login", controllers.AdminLogin)
	e.POST("/admin/logout", controllers.AdminLogout)
}
