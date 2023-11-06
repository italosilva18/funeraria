package models

type Plan struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	// Adicione outros campos relevantes do plano
}
