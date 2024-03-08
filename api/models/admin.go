package models

// Admin representa a estrutura de dados para um administrador do sistema
type Admin struct {
	ID    int    `json:"id"`
	Nome  string `json:"nome"`
	Email string `json:"email"`
	Senha string `json:"senha"`
	// Outros campos necessários
}
