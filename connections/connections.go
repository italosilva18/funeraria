package connections

import (
	"fmt"
	"margem/robo/connections/firebird"
	"margem/robo/connections/oracle"
	"margem/robo/connections/postgresql"
	"margem/robo/connections/sqlserver"
	"margem/robo/models/config"
)

// GetConnection retorna o driver de banco adequado para a automação e tipo de banco configurado.
func GetConnection(automacao string, loja config.Loja) (Database, error) {
	switch automacao {
	case "RMS":
		return oracle.New(loja), nil
	case "DEREVO":
		return firebird.New(loja), nil
	case "GESTOR":
		return firebird.New(loja), nil
	case "ARIUS":
		return postgresql.New(loja), nil
	case "ARIUS_ERP":
		return sqlserver.New(loja), nil
	default:
		return nil, fmt.Errorf("automação:'%s' não suportado para '%v'", automacao, loja)
	}
}
