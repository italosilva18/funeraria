package models

// Fatura representa a estrutura de dados para uma fatura de serviço prestado
type Fatura struct {
	ID          int     `json:"id"`
	ClienteID   int     `json:"cliente_id"`
	ServicoID   int     `json:"servico_id"`
	ValorTotal  float64 `json:"valor_total"`
	DataEmissao string  `json:"data_emissao"`
	Status      string  `json:"status"` // Pode ser "pendente", "pago", "atrasado", etc.
	// Outros campos necessários
}
