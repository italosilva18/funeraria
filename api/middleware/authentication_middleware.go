package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// AuthenticationMiddleware é um middleware para autenticar o administrador
func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Lógica de autenticação do administrador (implementação dependente do seu código)

		// Verifica se o token de acesso está presente nos cabeçalhos da requisição
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token de acesso não fornecido")
		}

		// Verifica se o token é válido (implementação dependente do seu código)
		if token != "seu_token_de_acesso" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token de acesso inválido")
		}

		// Se o token for válido, passa para o próximo middleware ou manipulador
		return next(c)
	}
}
