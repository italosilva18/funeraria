package oracle

import (
	"database/sql"
	"fmt"
	"margem/robo/models/config"

	go_ora "github.com/sijms/go-ora/v2"
)

type Oracle struct {
	loja config.Loja
}

// New cria uma instância de conexão Oracle com base nos dados da loja
func New(loja config.Loja) *Oracle {
	return &Oracle{loja: loja}
}

// DB abre uma conexão com banco Oracle utilizando os dados da loja
func (o *Oracle) DB() (*sql.DB, error) {
	connStr := go_ora.BuildUrl(
		o.loja.BancoFrenteLojaPrimario.Host,
		o.loja.BancoFrenteLojaPrimario.Porta,
		o.loja.BancoFrenteLojaPrimario.NomeBanco,
		o.loja.BancoFrenteLojaPrimario.Usuario,
		o.loja.BancoFrenteLojaPrimario.Senha,
		nil,
	)

	db, err := sql.Open("oracle", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão Oracle: %w", err)
	}

	// Testa a conexão imediatamente
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao validar conexão Oracle: %w", err)
	}

	return db, nil
}
