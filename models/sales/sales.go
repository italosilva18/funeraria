package sales

// Deduction - Modelo de desconto de item no MARGEM.
type Deduction struct {
	Description    string  `json:"descricao" bson:"descricao"`
	Ean            int64   `json:"codigo_ean" bson:"codigo_ean,omitempty"`
	Serial         int64   `json:"codigo_cupom" bson:"codigo_cupom"`
	Item           int64   `json:"item" bson:"item"`
	TotalDeduction float64 `json:"total_desconto" bson:"total_desconto"`
	Date           string  `json:"data" bson:"data"`
	Hour           string  `json:"hora" bson:"hora"`
	Operator       string  `json:"operador" bson:"operador"`
	Autorizer      *string `json:"supervisor" bson:"supervisor"`
}

// Salesman - Modelo de vendas por vendedor no MARGEM.
type Salesman struct {
	Name          string  `json:"nome" bson:"nome"`
	Amount        float64 `json:"total_vendido" bson:"total_vendido"`
	Coupons       int64   `json:"cupons" bson:"cupons"`
	Itens         int64   `json:"itens" bson:"itens"`
	AverageItems  float64 `json:"media_itens" bson:"media_itens"`
	AverageAmount float64 `json:"media_valor" bson:"media_valor"`
}
