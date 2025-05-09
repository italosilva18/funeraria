package rms

import (
	"fmt"
	"margem/robo/querys"
)

// GetQueryHandler retorna o handler de queries da versão correta do RMS com base no códigoSistema
func GetQueryHandler(codigoSistema string) querys.QueryHandler {
	switch codigoSistema {
	case "51.1":
		return &Rms511{}
	case "51.2":
		return &Rms512{}
	default:
		panic(fmt.Sprintf("❌ Versão %s não implementada para RMS.", codigoSistema))
	}
}
