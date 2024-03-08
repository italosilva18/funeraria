package models

// Agendamento representa a estrutura de dados para um agendamento de serviço funerário
type Agendamento struct {
	ID        int    `json:"id"`
	ClienteID int    `json:"cliente_id"`
	ServicoID int    `json:"servico_id"`
	Data      string `json:"data"`
	Hora      string `json:"hora"`
	Local     string `json:"local"`
	// Outros campos necessários
}
