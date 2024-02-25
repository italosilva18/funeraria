package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	// Configuração do template
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	e.Renderer = &TemplateRenderer{
		templates: tmpl,
	}

	// Rotas
	e.GET("/", homeHandler)
	e.GET("/cadastro", cadastroHandler)
	e.GET("/obito", obitoHandler)
	e.GET("/orcamento", orcamentoHandler)

	// Inicie o servidor
	e.Start(":8080")
}

// TemplateRenderer é uma estrutura para renderizar templates HTML
type TemplateRenderer struct {
	templates *template.Template
}

// Render renderiza o template HTML
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// Handlers
func homeHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "base.html", nil)
}

func cadastroHandler(c echo.Context) error {
	// Lógica para exibir a página de cadastro
	return c.Render(http.StatusOK, "cadastro.html", nil)
}

func obitoHandler(c echo.Context) error {
	// Lógica para exibir a página de registro de óbito
	return c.Render(http.StatusOK, "obito.html", nil)
}

func orcamentoHandler(c echo.Context) error {
	// Lógica para exibir a página de orçamento
	return c.Render(http.StatusOK, "orcamento.html", nil)
}
