package models

// Produto representa a estrutura de dados para um produto relacionado aos serviços funerários
type Produto struct {
	ID            int     `json:"id"`
	Nome          string  `json:"nome"`
	Descricao     string  `json:"descricao"`
	Preco         float64 `json:"preco"`
	Quantidade    int     `json:"quantidade"`
	UnidadeMedida string  `json:"unidade_medida"`
	// Outros campos necessários
}
