package models

type Sale struct {
	ID            int    `json:"id"`
	ProductID     int    `json:"product_id"`
	PlanID        int    `json:"plan_id"`
	CustomerID    int    `json:"customer_id"`
	TransactionID string `json:"transaction_id"`
	// Adicione outros campos relevantes da venda
}
