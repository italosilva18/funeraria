package section

// Sales - Modelos de vendas por seção no MARGEM.
type Sales struct {
	Code        interface{} `json:"codigo" bson:"codigo"`
	Description *string     `json:"descricao" bson:"descricao"`
	Quantity    int64       `json:"quantidade" bson:"quantidade"`
	Amout       float64     `json:"valor" bson:"valor"`
	Cost        float64     `json:"custo" bson:"custo,omitempty"`
	Profit      float64     `json:"lucro" bson:"lucro,omitempty"`
}
