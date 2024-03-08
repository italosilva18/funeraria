package controllers

import (
	"net/http"
	"strconv"

	"github.com/italosilva18/funeraria/api/models"
	"github.com/labstack/echo/v4"
)

// CriarServico cria um novo serviço
func CriarServico(c echo.Context) error {
	servico := new(models.Servico)
	if err := c.Bind(servico); err != nil {
		return err
	}

	// Salvar serviço no banco de dados (implementação dependente do seu código)

	return c.JSON(http.StatusCreated, servico)
}

// ObterServicos retorna todos os serviços
func ObterServicos(c echo.Context) error {
	// Obter todos os serviços do banco de dados (implementação dependente do seu código)

	servicos := []models.Servico{} // Supondo que você tenha uma função para obter todos os serviços
	return c.JSON(http.StatusOK, servicos)
}

// ObterServicoByID retorna um serviço específico por ID
func ObterServicoByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Obter serviço por ID do banco de dados (implementação dependente do seu código)

	servico := models.Servico{} // Supondo que você tenha uma função para obter um serviço por ID
	return c.JSON(http.StatusOK, servico)
}

// AtualizarServico atualiza os detalhes de um serviço existente
func AtualizarServico(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	servico := new(models.Servico)
	if err := c.Bind(servico); err != nil {
		return err
	}

	// Atualizar serviço no banco de dados (implementação dependente do seu código)

	return c.JSON(http.StatusOK, servico)
}

// DeletarServico exclui um serviço existente
func DeletarServico(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Excluir serviço do banco de dados (implementação dependente do seu código)

	return c.NoContent(http.StatusNoContent)
}
