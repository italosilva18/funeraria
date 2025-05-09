package entity

type CuponsPorOperador struct {
	QtdCupom float64 `json:"Qtd_Cupom" db:"qtd_cupom"`
	Operador string  `json:"Operador" db:"operador"`
}
