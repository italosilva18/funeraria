package derevo

import (
	"fmt"
	"margem/robo/querys"
)

// GetQueryHandler retorna o handler de queries da versão correta do RMS com base no códigoSistema
func GetQueryHandler(codigoSistema string) querys.QueryHandler {
	switch codigoSistema {
	case "52.1":
		return &Derevo521{}
	default:
		panic(fmt.Sprintf("❌ Versão %s não implementada para derevo.", codigoSistema))
	}
}
