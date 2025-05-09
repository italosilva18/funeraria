package product

import "time"

// Item - Modelo de produto no MARGEM.
type Item struct {
	Serial       string    `json:"licenca,omitempty" bson:"licenca,omitempty"`
	CNPJ         string    `json:"cnpj,omitempty" bson:"cnpj,omitempty"`
	State        string    `json:"estado,omitempty" bson:"estado,omitempty"`
	CodeState    int32     `json:"codigo_estado,omitempty" bson:"codigo_estado,omitempty"`
	City         string    `json:"cidade,omitempty" bson:"cidade,omitempty"`
	CodeCity     int32     `json:"codigo_cidade,omitempty" bson:"codigo_cidade,omitempty"`
	Neighborhood string    `json:"bairro,omitempty" bson:"bairro,omitempty"`
	Ean          int64     `json:"codigo_ean" bson:"codigo_ean"`
	Description  *string   `json:"descricao" bson:"descricao"`
	TotalAmount  float64   `json:"total_vendido" bson:"total_vendido"`
	Price        float64   `json:"preco_unitario" bson:"preco_unitario"`
	Quantity     float64   `json:"quantidade" bson:"quantidade"`
	Cost         float64   `json:"custo,omitempty" bson:"custo,omitempty"`
	Profit       float64   `json:"lucro,omitempty" bson:"lucro,omitempty"`
	Date         string    `json:"data,omitempty" bson:"data,omitempty"`
	Created      time.Time `json:"criado_em,omitempty" bson:"criado_em,omitempty"`
	Delete       time.Time `json:"expira,omitempty" bson:"expira,omitempty"`
}
