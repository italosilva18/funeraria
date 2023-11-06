package handlers

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

var secretKey = []byte("secretpassword") // Chave secreta para assinar e verificar tokens JWT

// CustomClaims define as informações que serão armazenadas em um token JWT
type CustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateToken gera um token JWT para um usuário autenticado
func GenerateToken(username string) (string, error) {
	claims := &CustomClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Token expira em 24 horas
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// AuthMiddleware é um middleware para autenticação
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token de autorização ausente"})
		}

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token inválido"})
		}

		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			// Autenticação bem-sucedida, você pode acessar as informações do usuário a partir de 'claims.Username'
			c.Set("username", claims.Username)
			return next(c)
		}

		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token inválido"})
	}
}
