package routes

import (
	"github.com/italosilva18/funeraria/api/controllers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	// Rota de cadastro de clientes
	e.POST("/clientes", controllers.CriarCliente)
	e.GET("/clientes", controllers.ObterClientes)
	e.GET("/clientes/:id", controllers.ObterClienteByID)
	e.PUT("/clientes/:id", controllers.AtualizarCliente)
	e.DELETE("/clientes/:id", controllers.DeletarCliente)

	// Rota de gestão de serviços funerários
	e.POST("/servicos", controllers.CriarServico)
	e.GET("/servicos", controllers.ObterServicos)
	e.GET("/servicos/:id", controllers.ObterServicoByID)
	e.PUT("/servicos/:id", controllers.AtualizarServico)
	e.DELETE("/servicos/:id", controllers.DeletarServico)

	// Rota de controle de estoque e produtos
	e.POST("/produtos", controllers.CriarProduto)
	e.GET("/produtos", controllers.ObterProdutos)
	e.GET("/produtos/:id", controllers.ObterProdutoByID)
	e.PUT("/produtos/:id", controllers.AtualizarProduto)
	e.DELETE("/produtos/:id", controllers.DeletarProduto)

	// Rota de agendamento e calendário
	e.POST("/agendamentos", controllers.CriarAgendamento)
	e.GET("/agendamentos", controllers.ObterAgendamentos)
	e.GET("/agendamentos/:id", controllers.ObterAgendamentoByID)
	e.PUT("/agendamentos/:id", controllers.AtualizarAgendamento)
	e.DELETE("/agendamentos/:id", controllers.DeletarAgendamento)

	// Rota de financeiro e cobranças
	e.POST("/faturas", controllers.CriarFatura)
	e.GET("/faturas", controllers.ObterFaturas)
	e.GET("/faturas/:id", controllers.ObterFaturaByID)
	e.PUT("/faturas/:id", controllers.AtualizarFatura)
	e.DELETE("/faturas/:id", controllers.DeletarFatura)

	// Rota de autenticação do administrador
	e.POST("/admin/login", controllers.AdminLogin)
	e.POST("/admin/logout", controllers.AdminLogout)

}
