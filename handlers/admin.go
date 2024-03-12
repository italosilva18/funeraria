package handlers

import (
	"net/http"

	"github.com/italosilva18/funeraria/models"

	"github.com/labstack/echo/v4"
)

// Função para cadastrar um admin
func CreateAdmin(c echo.Context) error {
	admin := new(models.Admin)
	if err := c.Bind(admin); err != nil {
		return err
	}
	// Aqui você adicionaria o admin ao banco de dados
	return c.JSON(http.StatusCreated, admin)
}
