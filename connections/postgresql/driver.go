package postgresql

import (
	"database/sql"
	"fmt"
	"margem/robo/models/config"

	_ "github.com/lib/pq"
)

type Postgres struct {
	loja config.Loja
}

func New(loja config.Loja) *Postgres {
	return &Postgres{loja: loja}
}

func (p *Postgres) DB() (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.loja.BancoFrenteLojaPrimario.Usuario,
		p.loja.BancoFrenteLojaPrimario.Senha,
		p.loja.BancoFrenteLojaPrimario.Host,
		p.loja.BancoFrenteLojaPrimario.Porta,
		p.loja.BancoFrenteLojaPrimario.NomeBanco,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão PostgreSQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao validar conexão PostgreSQL: %w", err)
	}

	return db, nil
}
