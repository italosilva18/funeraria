package arius_erp

import (
	"fmt"
	"margem/robo/querys"
)

// GetQueryHandler retorna o handler de queries da versão correta do RMS com base no códigoSistema
func GetQueryHandler(codigoSistema string) querys.QueryHandler {
	switch codigoSistema {
	case "1.3":
		return &AriusErp13{}

	default:
		panic(fmt.Sprintf("❌ Versão %s não implementada para ARIUS-ERP.", codigoSistema))
	}
}
