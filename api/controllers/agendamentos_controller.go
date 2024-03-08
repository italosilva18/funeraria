package controllers

import (
	"net/http"
	"strconv"

	"funeraria/api/models"

	"github.com/labstack/echo/v4"
)

// CriarAgendamento cria um novo agendamento
func CriarAgendamento(c echo.Context) error {
	agendamento := new(models.Agendamento)
	if err := c.Bind(agendamento); err != nil {
		return err
	}

	// Salvar agendamento no banco de dados (implementação dependente do seu código)

	return c.JSON(http.StatusCreated, agendamento)
}

// ObterAgendamentos retorna todos os agendamentos
func ObterAgendamentos(c echo.Context) error {
	// Obter todos os agendamentos do banco de dados (implementação dependente do seu código)

	agendamentos := []models.Agendamento{} // Supondo que você tenha uma função para obter todos os agendamentos
	return c.JSON(http.StatusOK, agendamentos)
}

// ObterAgendamentoByID retorna um agendamento específico por ID
func ObterAgendamentoByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Obter agendamento por ID do banco de dados (implementação dependente do seu código)

	agendamento := models.Agendamento{} // Supondo que você tenha uma função para obter um agendamento por ID
	return c.JSON(http.StatusOK, agendamento)
}

// AtualizarAgendamento atualiza os detalhes de um agendamento existente
func AtualizarAgendamento(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	agendamento := new(models.Agendamento)
	if err := c.Bind(agendamento); err != nil {
		return err
	}

	// Atualizar agendamento no banco de dados (implementação dependente do seu código)

	return c.JSON(http.StatusOK, agendamento)
}

// DeletarAgendamento exclui um agendamento existente
func DeletarAgendamento(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Excluir agendamento do banco de dados (implementação dependente do seu código)

	return c.NoContent(http.StatusNoContent)
}
