package controllers

import (
	"net/http"
	"strconv"

	"github.com/italosilva18/funeraria/api/models"

	"github.com/labstack/echo/v4"
)

// CriarFatura emite uma nova fatura para um serviço prestado
func CriarFatura(c echo.Context) error {
	fatura := new(models.Fatura)
	if err := c.Bind(fatura); err != nil {
		return err
	}

	// Salvar fatura no banco de dados (implementação dependente do seu código)

	return c.JSON(http.StatusCreated, fatura)
}

// ObterFaturas retorna todas as faturas
func ObterFaturas(c echo.Context) error {
	// Obter todas as faturas do banco de dados (implementação dependente do seu código)

	faturas := []models.Fatura{} // Supondo que você tenha uma função para obter todas as faturas
	return c.JSON(http.StatusOK, faturas)
}

// ObterFaturaByID retorna uma fatura específica por ID
func ObterFaturaByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Obter fatura por ID do banco de dados (implementação dependente do seu código)

	fatura := models.Fatura{} // Supondo que você tenha uma função para obter uma fatura por ID
	return c.JSON(http.StatusOK, fatura)
}

// AtualizarFatura atualiza os detalhes de uma fatura existente
func AtualizarFatura(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	fatura := new(models.Fatura)
	if err := c.Bind(fatura); err != nil {
		return err
	}

	// Atualizar fatura no banco de dados (implementação dependente do seu código)

	return c.JSON(http.StatusOK, fatura)
}

// DeletarFatura exclui uma fatura existente
func DeletarFatura(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Excluir fatura do banco de dados (implementação dependente do seu código)

	return c.NoContent(http.StatusNoContent)
}
