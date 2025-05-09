package hour

// Sales - Modelo de vendas por Hora do MARGEM.
type Sales struct {
	Hour          int64   `json:"hora" bson:"hora"`
	Amount        float64 `json:"total_vendido" bson:"total_vendido"`
	Cupons        int64   `json:"total_cupons" bson:"total_cupons"`
	Itens         int64   `json:"total_itens" bson:"total_itens"`
	AverageItens  float64 `json:"media_itens" bson:"media_itens"`
	AverageAmount float64 `json:"media_valor" bson:"media_valor"`
}
