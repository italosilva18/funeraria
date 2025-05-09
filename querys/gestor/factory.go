package gestor

import (
	"fmt"
	"margem/robo/querys"
)

// GetQueryHandler retorna o handler de queries da versão correta do RMS com base no códigoSistema
func GetQueryHandler(codigoSistema string) querys.QueryHandler {
	switch codigoSistema {
	case "48.1":
		return &Gestor481{}
	case "48.2":
		return &Gestor482{}
	default:
		panic(fmt.Sprintf("❌ Versão %s não implementada para GESTOR.", codigoSistema))
	}
}
