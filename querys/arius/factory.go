package arius

import (
	"fmt"
	"margem/robo/querys"
)

// GetQueryHandler retorna o handler de queries da versão correta do RMS com base no códigoSistema
func GetQueryHandler(codigoSistema string) querys.QueryHandler {
	switch codigoSistema {
	case "1.1":
		return &Arius11{}
	case "1.2":
		return &Arius12{}
	default:
		panic(fmt.Sprintf("❌ Versão %s não implementada para ARIUS.", codigoSistema))
	}
}
