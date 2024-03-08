package models

// Servico representa a estrutura de dados para um serviço funerário
type Servico struct {
	ID          int     `json:"id"`
	Tipo        string  `json:"tipo"`
	Descricao   string  `json:"descricao"`
	Preco       float64 `json:"preco"`
	Responsavel string  `json:"responsavel"`
	// Outros campos necessários
}
