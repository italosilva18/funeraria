package models

// Cliente representa a estrutura de dados para um cliente
type Cliente struct {
	ID             int    `json:"id"`
	NomeCompleto   string `json:"nome_completo"`
	CPF            string `json:"cpf"`
	DataNascimento string `json:"data_nascimento"`
	Endereco       string `json:"endereco"`
	Telefone       string `json:"telefone"`
	Email          string `json:"email"`
	// Outros campos necessários
}
