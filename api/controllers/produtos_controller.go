package controllers

import (
	"net/http"
	"strconv"

	"github.com/italosilva18/funeraria/api/models"
	"github.com/labstack/echo/v4"
)

// CriarProduto cria um novo produto
func CriarProduto(c echo.Context) error {
	produto := new(models.Produto)
	if err := c.Bind(produto); err != nil {
		return err
	}

	// Salvar produto no banco de dados (implementação dependente do seu código)

	return c.JSON(http.StatusCreated, produto)
}

// ObterProdutos retorna todos os produtos
func ObterProdutos(c echo.Context) error {
	// Obter todos os produtos do banco de dados (implementação dependente do seu código)

	produtos := []models.Produto{} // Supondo que você tenha uma função para obter todos os produtos
	return c.JSON(http.StatusOK, produtos)
}

// ObterProdutoByID retorna um produto específico por ID
func ObterProdutoByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Obter produto por ID do banco de dados (implementação dependente do seu código)

	produto := models.Produto{} // Supondo que você tenha uma função para obter um produto por ID
	return c.JSON(http.StatusOK, produto)
}

// AtualizarProduto atualiza os detalhes de um produto existente
func AtualizarProduto(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	produto := new(models.Produto)
	if err := c.Bind(produto); err != nil {
		return err
	}

	// Atualizar produto no banco de dados (implementação dependente do seu código)

	return c.JSON(http.StatusOK, produto)
}

// DeletarProduto exclui um produto existente
func DeletarProduto(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Excluir produto do banco de dados (implementação dependente do seu código)

	return c.NoContent(http.StatusNoContent)
}
