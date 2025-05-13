package concinco

import (
	"fmt"
	"margem/robo/querys"
)

// GetQueryHandler retorna o handler de queries da versão correta do CONCINCO com base no códigoSistema
func GetQueryHandler(codigoSistema string) querys.QueryHandler {
	switch codigoSistema {
	case "28.1":
		return &Concinco281{}
	default:
		panic(fmt.Sprintf("❌ Versão %s não implementada para CONCINCO.", codigoSistema))
	}
}
