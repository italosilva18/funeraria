package firebird

import (
	"database/sql"
	"fmt"
	"margem/robo/models/config"

	_ "github.com/nakagami/firebirdsql"
)

type Firebird struct {
	loja config.Loja
}

func New(loja config.Loja) *Firebird {
	return &Firebird{loja: loja}
}

func (f *Firebird) DB() (*sql.DB, error) {
	// Constrói a string de conexão incluindo a porta e o charset UTF8
	connStr := fmt.Sprintf(
		"%s:%s@%s:%d/%s?charset=UTF8",
		f.loja.BancoFrenteLojaPrimario.Usuario,
		f.loja.BancoFrenteLojaPrimario.Senha,
		f.loja.BancoFrenteLojaPrimario.Host,
		f.loja.BancoFrenteLojaPrimario.Porta,
		f.loja.BancoFrenteLojaPrimario.NomeBanco,
	)

	db, err := sql.Open("firebirdsql", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão Firebird: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("erro ao validar conexão Firebird: %w", err)
	}

	return db, nil
}
