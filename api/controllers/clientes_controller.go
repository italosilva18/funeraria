package controllers

import (
	"net/http"
	"strconv"

	"github.com/italosilva18/funeraria/api/models"
	"github.com/labstack/echo/v4"
)

// CriarCliente cria um novo cliente
func CriarCliente(c echo.Context) error {
	cliente := new(models.Cliente)
	if err := c.Bind(cliente); err != nil {
		return err
	}

	// Salvar cliente no banco de dados (implementação dependente do seu código)

	return c.JSON(http.StatusCreated, cliente)
}

// ObterClientes retorna todos os clientes
func ObterClientes(c echo.Context) error {
	// Obter todos os clientes do banco de dados (implementação dependente do seu código)

	clientes := []models.Cliente{} // Supondo que você tenha uma função para obter todos os clientes
	return c.JSON(http.StatusOK, clientes)
}

// ObterClienteByID retorna um cliente específico por ID
func ObterClienteByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Obter cliente por ID do banco de dados (implementação dependente do seu código)

	cliente := models.Cliente{} // Supondo que você tenha uma função para obter um cliente por ID
	return c.JSON(http.StatusOK, cliente)
}

// AtualizarCliente atualiza os detalhes de um cliente existente
func AtualizarCliente(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	cliente := new(models.Cliente)
	if err := c.Bind(cliente); err != nil {
		return err
	}

	// Atualizar cliente no banco de dados (implementação dependente do seu código)

	return c.JSON(http.StatusOK, cliente)
}

// DeletarCliente exclui um cliente existente
func DeletarCliente(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Excluir cliente do banco de dados (implementação dependente do seu código)

	return c.NoContent(http.StatusNoContent)
}
