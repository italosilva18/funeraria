package coupon

// Operator -
type Operator struct {
	Name     string `json:"nome_operador" bson:"nome_operador"`
	Quantity int64  `json:"quantidade" bson:"quantidade"`
}
