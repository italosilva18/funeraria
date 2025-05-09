package mysql

import (
	"database/sql"
	"fmt"
	"margem/robo/models/config"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	loja config.Loja
}

func New(loja config.Loja) *MySQL {
	return &MySQL{loja: loja}
}

func (m *MySQL) DB() (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		m.loja.BancoFrenteLojaPrimario.Usuario,
		m.loja.BancoFrenteLojaPrimario.Senha,
		m.loja.BancoFrenteLojaPrimario.Host,
		m.loja.BancoFrenteLojaPrimario.Porta,
		m.loja.BancoFrenteLojaPrimario.NomeBanco,
	)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão MySQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao validar conexão MySQL: %w", err)
	}

	return db, nil
}
