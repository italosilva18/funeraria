package models

import (
	"margem/robo/models/coupon"
	"margem/robo/models/hour"
	"margem/robo/models/modality"
	"margem/robo/models/product"
	"margem/robo/models/reversed"
	"margem/robo/models/sales"
	"margem/robo/models/section"
	"time"
)

// Item - Representa um relatório no banco de dados.
type Item struct {
	Store              string              `json:"loja" bson:"loja"`
	Serial             string              `json:"licenca" bson:"licenca"`
	CNPJ               string              `json:"cnpj" bson:"cnpj"`
	Date               string              `json:"data" bson:"data"`
	State              string              `json:"estado" bson:"estado"`
	CodeState          int32               `json:"codigo_estado,omitempty" bson:"codigo_estado"`
	City               string              `json:"cidade" bson:"cidade"`
	CodeCity           int32               `json:"codigo_cidade,omitempty" bson:"codigo_cidade"`
	Neighborhood       string              `json:"bairro" bson:"bairro"`
	Partner            string              `json:"automacao" bson:"automacao"`
	Version            string              `json:"versao" bson:"versao"`
	TotalAmount        float64             `json:"total_vendido" bson:"total_vendido"`
	Cost               float64             `json:"custo" bson:"custo"`
	Profit             float64             `json:"lucro" bson:"lucro"`
	Clients            int64               `json:"clientes" bson:"clientes"`
	AverageTicket      float64             `json:"ticket_medio" bson:"ticket_medio"`
	ProductsService    float64             `json:"produto_atendimento" bson:"produto_atendimento"`
	CanceledCoupons    int64               `json:"cupons_cancelados" bson:"cupons_cancelados"`
	HourRange          []hour.Sales        `json:"faixa_horaria" bson:"faixa_horaria"`
	Modalities         []modality.Sales    `json:"modalidades_pagamento" bson:"modalidades_pagamento"`
	ModalitiesOperator []modality.Operator `json:"modalidades_operador" bson:"modalidades_operador"`
	ReversedCoupons    []reversed.Coupons  `json:"estornos_cupons" bson:"estornos_cupons"`
	Deduction          []sales.Deduction   `json:"descontos" bson:"descontos"`
	CouponsOperator    []coupon.Operator   `json:"cupons_operador" bson:"cupons_operador"`
	Sections           []section.Sales     `json:"sessao" bson:"sessao"`
	Salesman           []sales.Salesman    `json:"vendedores" bson:"vendedores"`
	Products           []product.Item      `json:"rank_produtos" bson:"rank_produtos"`
	CreateAt           time.Time           `json:"criado_em" bson:"criado_em"`
	Delete             time.Time           `json:"auto_delete" bson:"auto_delete"`
	Synce              time.Time           `json:"sincronizado,omitempty"`
}
