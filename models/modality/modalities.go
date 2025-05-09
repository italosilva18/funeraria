package modality

// Sales - Modelo de vendas por modalidade no MARGEM.
type Sales struct {
	Code        interface{} `json:"codigo_modalidade" bson:"codigo_modalidade"`
	Description string      `json:"descricao" bson:"descricao"`
	Quantity    int64       `json:"quantidade" bson:"quantidade"`
	TotalAmount float64     `json:"total_vendido" bson:"total_vendido"`
	Cost        float64     `json:"custo" bson:"custo,omitempty"`
	Profit      float64     `json:"lucro" bson:"lucro,omitempty"`
}

// Operator - Modelo de vendas por operador de caixa no MARGEM.
type Operator struct {
	Operator    string      `json:"operador" bson:"operador"`
	Code        interface{} `json:"codigo_modalidade" bson:"codigo_modalidade"`
	Description string      `json:"descricao" bson:"descricao"`
	Quantity    int64       `json:"quantidade" bson:"quantidade"`
	TotalAmount float64     `json:"total_vendido" bson:"total_vendido"`
}
