package sqlserver

import (
	"database/sql"
	"fmt"
	"margem/robo/models/config"

	_ "github.com/microsoft/go-mssqldb"
)

type SQLServer struct {
	loja config.Loja
}

func New(loja config.Loja) *SQLServer {
	return &SQLServer{loja: loja}
}

func (s *SQLServer) DB() (*sql.DB, error) {
	connStr := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		s.loja.BancoFrenteLojaPrimario.Usuario,
		s.loja.BancoFrenteLojaPrimario.Senha,
		s.loja.BancoFrenteLojaPrimario.Host,
		s.loja.BancoFrenteLojaPrimario.Porta,
		s.loja.BancoFrenteLojaPrimario.NomeBanco,
	)

	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão SQL Server: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao validar conexão SQL Server: %w", err)
	}

	return db, nil
}
