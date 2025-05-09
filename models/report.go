package models

import "time"

type Report struct {
	Data                 string                 `json:"data"`
	Licenca              string                 `json:"licenca"`
	AutoDelete           time.Time              `json:"auto_delete"`
	Automacao            string                 `json:"automacao"`
	Bairro               string                 `json:"bairro"`
	Cidade               string                 `json:"cidade"`
	Clientes             int                    `json:"clientes"`
	Cnpj                 string                 `json:"cnpj"`
	CodigoCidade         int                    `json:"codigo_cidade"`
	CodigoEstado         int                    `json:"codigo_estado"`
	CriadoEm             time.Time              `json:"criado_em"`
	CuponsCancelados     float64                `json:"cupons_cancelados"`
	CuponsOperador       []CuponsOperador       `json:"cupons_operador"`
	Custo                float64                `json:"custo"`
	Estado               string                 `json:"estado"`
	EstornosCupons       EstornosCupons         `json:"estornos_cupons"`
	FaixaHoraria         FaixaHoraria           `json:"faixa_horaria"`
	Loja                 string                 `json:"loja"`
	Lucro                float64                `json:"lucro"`
	ModalidadesOperador  []ModalidadesOperador  `json:"modalidades_operador"`
	ModalidadesPagamento []ModalidadesPagamento `json:"modalidades_pagamento"`
	ProdutoAtendimento   float64                `json:"produto_atendimento"`
	RankProdutos         []RankProdutos         `json:"rank_produtos"`
	Sessao               []Sessao               `json:"sessao"`
}

type CuponsOperador struct {
	NomeOperador string  `json:"nome_operador"`
	Quantidade   float64 `json:"quantidade"`
}

type Descontos struct {
}

type EstornosCupons struct {
	Pdv         string `json:"pdv"`
	Operador    string `json:"operador"`
	Autorizador string `json:"autorizador"`
	CodigoCupom string `json:"codigo_cupom"`
	Data        string `json:"data"`
	Hora        string `json:"hora"`
	Valor       string `json:"valor"`
}

type FaixaHoraria struct {
	TotalCupons  float64 `json:"total_cupons"`
	TotalItens   float64 `json:"total_itens"`
	MediaItens   float64 `json:"media_itens"`
	MediaValor   float64 `json:"media_valor"`
	Hora         float64 `json:"hora"`
	TotalVendido float64 `json:"total_vendido"`
}

type ModalidadesOperador struct {
	TotalVendido     float64 `json:"total_vendido"`
	Operador         string  `json:"operador"`
	CodigoModalidade float64 `json:"codigo_modalidade"`
	Descricao        string  `json:"descricao"`
	Quantidade       float64 `json:"quantidade"`
}

type ModalidadesPagamento struct {
	Custo            float64 `json:"custo"`
	Lucro            float64 `json:"lucro"`
	CodigoModalidade float64 `json:"codigo_modalidade"`
	Descricao        string  `json:"descricao"`
	Quantidade       float64 `json:"quantidade"`
	TotalVendido     float64 `json:"total_vendido"`
}

type RankProdutos struct {
	Expira        string  `json:"expira"`
	CodigoEan     float64 `json:"codigo_ean"`
	Descricao     string  `json:"descricao"`
	TotalVendido  float64 `json:"total_vendido"`
	PrecoUnitario float64 `json:"preco_unitario"`
	Quantidade    float64 `json:"quantidade"`
	CriadoEm      string  `json:"criado_em"`
}

type Sessao struct {
	Custo      float64 `json:"custo"`
	Lucro      float64 `json:"lucro"`
	Codigo     float64 `json:"codigo"`
	Descricao  string  `json:"descricao"`
	Quantidade float64 `json:"quantidade"`
	Valor      float64 `json:"valor"`
}
