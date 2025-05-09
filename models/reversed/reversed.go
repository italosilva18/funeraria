package reversed

// Coupons -
type Coupons struct {
	Pdv        string  `json:"pdv" bson:"pdv"`
	Operator   string  `json:"operador" bson:"operador"`
	Authorizer *string `json:"autorizador" bson:"autorizador"`
	Code       string  `json:"codigo_cupom" bson:"codigo_cupom"`
	Date       string  `json:"data" bson:"data"`
	Hour       string  `json:"hora" bson:"hora"`
	Amount     float64 `json:"valor" bson:"valor"`
}
