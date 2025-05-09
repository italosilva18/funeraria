package repository

import (
	"database/sql"
	"fmt"
	"log"
	"margem/robo/models"
	"margem/robo/models/coupon"
	"margem/robo/models/hour"
	"margem/robo/models/modality"
	"margem/robo/models/product"
	"margem/robo/models/reversed"
	"margem/robo/models/sales"
	"margem/robo/models/section"
	"margem/robo/querys"
	"time"

	"github.com/shopspring/decimal"
)

// reportDerevoRepository implementa o ReportRepository para o sistema Derevo (versão 521).
type reportDerevoRepository struct {
	db    *sql.DB
	query querys.QueryHandler
}

// NewReportDerevoRepository retorna uma instância do repositório Derevo.
func NewReportDerevoRepository(db *sql.DB, q querys.QueryHandler) ReportRepository {
	return &reportDerevoRepository{
		db:    db,
		query: q,
	}
}

// FindAllRangeDay executa as 15 queries para montar o relatório (models.Item).
func (r *reportDerevoRepository) FindAllRangeDay(startDate, endDate, store string, _ *sql.DB) (models.Item, error) {
	// Verifica a conexão com o banco
	if r.db == nil {
		log.Println("❌ Erro: conexão com o banco de dados é nil")
		return models.Item{}, fmt.Errorf("conexão com o banco de dados é nil")
	}

	// Verifica o handler de queries
	if r.query == nil {
		log.Println("❌ Erro: handler de queries é nil")
		return models.Item{}, fmt.Errorf("handler de queries é nil")
	}

	// Estrutura inicial do relatório
	report := models.Item{
		Cost:            0,
		Profit:          0,
		TotalAmount:     0,
		Clients:         0,
		CanceledCoupons: 0,
		AverageTicket:   0,
		ProductsService: 0,
	}

	// Normaliza as datas
	startSimple, endSimple, err := r.normalizeDates(startDate, endDate)
	if err != nil {
		log.Printf("❌ Erro ao normalizar datas: %v", err)
		return report, err
	}

	// (1) TOTAL_VENDIDO_DIA
	lucro, custo, valor, err := r.lucroCustoDia(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar lucroCustoDia: %v", err)
		return report, err
	}
	report.Profit = lucro
	report.Cost = custo
	report.TotalAmount = valor

	// (2) TICKET_MEDIO_DIA
	tk, err := r.ticketMedio(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar ticketMedio: %v", err)
		return report, err
	}
	report.AverageTicket = tk

	// (3) PRODUTOS_POR_ATENDIMENTO_DIA
	ppa, err := r.produtoAtendimento(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar produtoAtendimento: %v", err)
		return report, err
	}
	report.ProductsService = ppa

	// (4) TOTAL_CUPONS_VALIDOS_DIA
	validCupom, err := r.totalCupomValido(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar totalCupomValido: %v", err)
		return report, err
	}
	report.Clients += int64(validCupom)

	// (5) TOTAL_CUPONS_CANCELADOS_DIA
	canc, err := r.cuponsCancelados(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar cuponsCancelados: %v", err)
		return report, err
	}
	report.CanceledCoupons = int64(canc)

	// (6) VENDAS_POR_FAIXA_HORARIA
	faixa, err := r.faixaHoraria(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar faixaHoraria: %v", err)
		return report, err
	}
	report.HourRange = *faixa

	// (7) PRODUTOS_MAIS_VENDIDO_DIA
	rank, err := r.rankProdutos(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar rankProdutos: %v", err)
		return report, err
	}
	report.Products = *rank

	// (8) PRODUTOS_ESTATISTICA
	_, err = r.produtosEstatistica(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar produtosEstatistica: %v", err)
		// Continua mesmo com erro, pois é opcional
	}

	// (9) TOTAL_VENDIDO_POR_MODALIDADES
	mods, err := r.modalidadesPagamento(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar modalidadesPagamento: %v", err)
		return report, err
	}
	report.Modalities = *mods

	// (10) ESTORNO_DE_CUPONS_POR_OPERADOR_DIA
	estornos, err := r.estornoCupons(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar estornoCupons: %v", err)
		return report, err
	}
	report.ReversedCoupons = *estornos

	// (11) DESCONTOS_DE_ITENS_DIA
	descontos, err := r.descontoItem(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar descontoItem: %v", err)
		return report, err
	}
	report.Deduction = *descontos

	// (12) VENDAS_POR_SECOES_DIA
	secs, err := r.section(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar section: %v", err)
		return report, err
	}
	report.Sections = *secs

	// (13) TOTAL_CUPONS_POR_OPERADOR_DIA
	cuponsOp, err := r.cuponsPorOperador(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar cuponsPorOperador: %v", err)
		return report, err
	}
	report.CouponsOperator = *cuponsOp
	for _, op := range *cuponsOp {
		report.Clients += op.Quantity
	}

	// (14) VENDAS_POR_VENDEDOR_DIA
	vend, err := r.vendedores(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar vendedores: %v", err)
		return report, err
	}
	report.Salesman = *vend

	// (15) MODALIDADES_PAGAMENTO_OPERADOR_DIA
	modOper, err := r.modalidadeOperador(startSimple, endSimple, store)
	if err != nil {
		log.Printf("❌ Erro ao executar modalidadeOperador: %v", err)
		return report, err
	}
	report.ModalitiesOperator = *modOper

	return report, nil
}

// normalizeDates converte datas para o formato YYYY-MM-DD
func (r *reportDerevoRepository) normalizeDates(s, e string) (string, string, error) {
	layoutFull := "2006-01-02 15:04:05"
	layoutShort := "2006-01-02"

	t1, err1 := time.Parse(layoutFull, s)
	t2, err2 := time.Parse(layoutFull, e)
	if err1 == nil && err2 == nil {
		return t1.Format(layoutShort), t2.Format(layoutShort), nil
	}

	_, err3 := time.Parse(layoutShort, s)
	_, err4 := time.Parse(layoutShort, e)
	if err3 == nil && err4 == nil {
		return s, e, nil
	}

	return "", "", fmt.Errorf("datas em formato inválido: '%s', '%s'", s, e)
}

// (1) TOTAL_VENDIDO_DIA
func (r *reportDerevoRepository) lucroCustoDia(start, end, store string) (float64, float64, float64, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_VENDIDO_DIA())
	if err != nil {
		return 0, 0, 0, fmt.Errorf("erro ao preparar query TOTAL_VENDIDO_DIA: %w", err)
	}
	defer stmt.Close()

	var data struct {
		Valor float64
		Custo float64
		Lucro float64
	}
	err = stmt.QueryRow(store, start, end).Scan(&data.Valor, &data.Custo, &data.Lucro)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("erro ao executar query TOTAL_VENDIDO_DIA: %w", err)
	}
	return data.Lucro, data.Custo, data.Valor, nil
}

// (2) TICKET_MEDIO_DIA
func (r *reportDerevoRepository) ticketMedio(start, end, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.TICKET_MEDIO_DIA())
	if err != nil {
		return 0, fmt.Errorf("erro ao preparar query TICKET_MEDIO_DIA: %w", err)
	}
	defer stmt.Close()

	var value float64
	err = stmt.QueryRow(store, start, end).Scan(&value)
	if err != nil {
		return 0, fmt.Errorf("erro ao executar query TICKET_MEDIO_DIA: %w", err)
	}
	return value, nil
}

// (3) PRODUTOS_POR_ATENDIMENTO_DIA
func (r *reportDerevoRepository) produtoAtendimento(start, end, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.PRODUTOS_POR_ATENDIMENTO_DIA())
	if err != nil {
		return 0, fmt.Errorf("erro ao preparar query PRODUTOS_POR_ATENDIMENTO_DIA: %w", err)
	}
	defer stmt.Close()

	var val float64
	err = stmt.QueryRow(store, start, end).Scan(&val)
	if err != nil {
		return 0, fmt.Errorf("erro ao executar query PRODUTOS_POR_ATENDIMENTO_DIA: %w", err)
	}
	return val, nil
}

// (4) TOTAL_CUPONS_VALIDOS_DIA
func (r *reportDerevoRepository) totalCupomValido(start, end, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_CUPONS_VALIDOS_DIA())
	if err != nil {
		return 0, fmt.Errorf("erro ao preparar query TOTAL_CUPONS_VALIDOS_DIA: %w", err)
	}
	defer stmt.Close()

	var val float64
	err = stmt.QueryRow(store, start, end).Scan(&val)
	if err != nil {
		return 0, fmt.Errorf("erro ao executar query TOTAL_CUPONS_VALIDOS_DIA: %w", err)
	}
	return val, nil
}

// (5) TOTAL_CUPONS_CANCELADOS_DIA
func (r *reportDerevoRepository) cuponsCancelados(start, end, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_CUPONS_CANCELADOS_DIA())
	if err != nil {
		return 0, fmt.Errorf("erro ao preparar query TOTAL_CUPONS_CANCELADOS_DIA: %w", err)
	}
	defer stmt.Close()

	var val float64
	err = stmt.QueryRow(store, start, end).Scan(&val)
	if err != nil {
		return 0, fmt.Errorf("erro ao executar query TOTAL_CUPONS_CANCELADOS_DIA: %w", err)
	}
	return val, nil
}

// (6) VENDAS_POR_FAIXA_HORARIA
func (r *reportDerevoRepository) faixaHoraria(start, end, store string) (*[]hour.Sales, error) {
	stmt, err := r.db.Prepare(r.query.VENDAS_POR_FAIXA_HORARIA())
	if err != nil {
		return nil, fmt.Errorf("erro ao preparar query VENDAS_POR_FAIXA_HORARIA: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(store, start, end)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query VENDAS_POR_FAIXA_HORARIA: %w", err)
	}
	defer rows.Close()

	var data []hour.Sales
	for rows.Next() {
		var e hour.Sales
		if err := rows.Scan(&e.Hour, &e.Amount, &e.Cupons, &e.Itens, &e.AverageItens, &e.AverageAmount); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de VENDAS_POR_FAIXA_HORARIA: %w", err)
		}
		data = append(data, e)
	}
	return &data, nil
}

// (7) PRODUTOS_MAIS_VENDIDO_DIA
func (r *reportDerevoRepository) rankProdutos(start, end, store string) (*[]product.Item, error) {
	stmt, err := r.db.Prepare(r.query.PRODUTOS_MAIS_VENDIDO_DIA())
	if err != nil {
		return nil, fmt.Errorf("erro ao preparar query PRODUTOS_MAIS_VENDIDO_DIA: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(store, start, end)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query PRODUTOS_MAIS_VENDIDO_DIA: %w", err)
	}
	defer rows.Close()

	var data []product.Item
	for rows.Next() {
		var e product.Item
		if err := rows.Scan(&e.Ean, &e.Description, &e.Price, &e.Quantity, &e.TotalAmount, &e.Cost, &e.Profit); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de PRODUTOS_MAIS_VENDIDO_DIA: %w", err)
		}
		data = append(data, e)
	}
	return &data, nil
}

// (8) PRODUTOS_ESTATISTICA
func (r *reportDerevoRepository) produtosEstatistica(start, end, store string) (*[]product.Item, error) {
	stmt, err := r.db.Prepare(r.query.PRODUTOS_ESTATISTICA())
	if err != nil {
		return nil, fmt.Errorf("erro ao preparar query PRODUTOS_ESTATISTICA: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(store, start, end)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query PRODUTOS_ESTATISTICA: %w", err)
	}
	defer rows.Close()

	var data []product.Item
	for rows.Next() {
		var e product.Item
		if err := rows.Scan(&e.Ean, &e.Description, &e.Price, &e.Quantity, &e.TotalAmount, &e.Cost, &e.Profit); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de PRODUTOS_ESTATISTICA: %w", err)
		}
		data = append(data, e)
	}
	return &data, nil
}

// (9) TOTAL_VENDIDO_POR_MODALIDADES
func (r *reportDerevoRepository) modalidadesPagamento(start, end, store string) (*[]modality.Sales, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_VENDIDO_POR_MODALIDADES())
	if err != nil {
		return nil, fmt.Errorf("erro ao preparar query TOTAL_VENDIDO_POR_MODALIDADES: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(store, start, end)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query TOTAL_VENDIDO_POR_MODALIDADES: %w", err)
	}
	defer rows.Close()

	var data []modality.Sales
	for rows.Next() {
		var e modality.Sales
		if err := rows.Scan(&e.Code, &e.Description, &e.Quantity, &e.TotalAmount); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de TOTAL_VENDIDO_POR_MODALIDADES: %w", err)
		}
		data = append(data, e)
	}
	return &data, nil
}

// (10) ESTORNO_DE_CUPONS_POR_OPERADOR_DIA
func (r *reportDerevoRepository) estornoCupons(start, end, store string) (*[]reversed.Coupons, error) {
	stmt, err := r.db.Prepare(r.query.ESTORNO_DE_CUPONS_POR_OPERADOR_DIA())
	if err != nil {
		return nil, fmt.Errorf("erro ao preparar query ESTORNO_DE_CUPONS_POR_OPERADOR_DIA: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(store, start, end)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query ESTORNO_DE_CUPONS_POR_OPERADOR_DIA: %w", err)
	}
	defer rows.Close()

	var data []reversed.Coupons
	for rows.Next() {
		var e reversed.Coupons
		var authorizer *string
		if err := rows.Scan(&e.Pdv, &e.Date, &e.Hour, &e.Operator, &authorizer, &e.Amount, &e.Code); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de ESTORNO_DE_CUPONS_POR_OPERADOR_DIA: %w", err)
		}
		e.Authorizer = authorizer
		data = append(data, e)
	}
	return &data, nil
}

// (11) DESCONTOS_DE_ITENS_DIA
func (r *reportDerevoRepository) descontoItem(start, end, store string) (*[]sales.Deduction, error) {
	stmt, err := r.db.Prepare(r.query.DESCONTOS_DE_ITENS_DIA())
	if err != nil {
		return nil, fmt.Errorf("erro ao preparar query DESCONTOS_DE_ITENS_DIA: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(store, start, end)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query DESCONTOS_DE_ITENS_DIA: %w", err)
	}
	defer rows.Close()

	var data []sales.Deduction
	for rows.Next() {
		var e sales.Deduction
		var authorizer *string
		if err := rows.Scan(&e.Date, &e.Hour, &e.Operator, &authorizer, &e.TotalDeduction, &e.Serial, &e.Item, &e.Ean, &e.Description); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de DESCONTOS_DE_ITENS_DIA: %w", err)
		}
		e.Autorizer = authorizer
		data = append(data, e)
	}
	return &data, nil
}

// (12) VENDAS POR SEÇÕES DIA
func (r *reportDerevoRepository) section(start, end, store string) (*[]section.Sales, error) {
	stmt, err := r.db.Prepare(r.query.VENDAS_POR_SECOES_DIA())
	if err != nil {
		return nil, fmt.Errorf("erro ao preparar query VENDAS_POR_SECOES_DIA: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(store, start, end)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query VENDAS_POR_SECOES_DIA: %w", err)
	}
	defer rows.Close()

	var data []section.Sales
	for rows.Next() {
		var e section.Sales
		var qtde float64
		if err := rows.Scan(&e.Code, &e.Description, &qtde, &e.Amout, &e.Cost, &e.Profit); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de VENDAS_POR_SECOES_DIA: %w", err)
		}
		e.Quantity = decimal.NewFromFloat(qtde).IntPart()
		data = append(data, e)
	}
	return &data, nil
}

// (13) TOTAL CUPONS POR OPERADOR DIA
func (r *reportDerevoRepository) cuponsPorOperador(start, end, store string) (*[]coupon.Operator, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_CUPONS_POR_OPERADOR_DIA())
	if err != nil {
		return nil, fmt.Errorf("erro ao preparar query TOTAL_CUPONS_POR_OPERADOR_DIA: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(store, start, end)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query TOTAL_CUPONS_POR_OPERADOR_DIA: %w", err)
	}
	defer rows.Close()

	var data []coupon.Operator
	for rows.Next() {
		var e coupon.Operator
		var dummy float64 // Ignora o campo 'pvalor'
		if err := rows.Scan(&e.Quantity, &e.Name, &dummy); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de TOTAL_CUPONS_POR_OPERADOR_DIA: %w", err)
		}
		data = append(data, e)
	}
	return &data, nil
}

// (14) VENDAS POR VENDEDOR DIA
func (r *reportDerevoRepository) vendedores(start, end, store string) (*[]sales.Salesman, error) {
	// Verifica se a query para vendedores está definida; para Derevo 521, espera-se que ela seja vazia.
	query := r.query.VENDAS_POR_VENDEDOR_DIA()
	if query == "" {
		return &[]sales.Salesman{}, nil
	}
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("erro ao preparar query VENDAS_POR_VENDEDOR_DIA: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(store, start, end, start, end)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query VENDAS_POR_VENDEDOR_DIA: %w", err)
	}
	defer rows.Close()

	var data []sales.Salesman
	var dataRegister string
	for rows.Next() {
		var e sales.Salesman
		if err := rows.Scan(&dataRegister, &e.Name, &e.Amount, &e.Coupons, &e.Itens, &e.AverageItems, &e.AverageAmount); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de VENDAS_POR_VENDEDOR_DIA: %w", err)
		}
		data = append(data, e)
	}
	return &data, nil
}

// (15) MODALIDADES DE PAGAMENTO POR OPERADOR DIA
func (r *reportDerevoRepository) modalidadeOperador(start, end, store string) (*[]modality.Operator, error) {
	stmt, err := r.db.Prepare(r.query.MODALIDADES_PAGAMENTO_OPERADOR_DIA())
	if err != nil {
		return nil, fmt.Errorf("erro ao preparar query MODALIDADES_PAGAMENTO_OPERADOR_DIA: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(store, start, end)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query MODALIDADES_PAGAMENTO_OPERADOR_DIA: %w", err)
	}
	defer rows.Close()

	var data []modality.Operator
	for rows.Next() {
		var e modality.Operator
		var dataStr string // Ignora o campo 'pdata'
		if err := rows.Scan(&dataStr, &e.Operator, &e.Description, &e.Quantity, &e.TotalAmount); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de MODALIDADES_PAGAMENTO_OPERADOR_DIA: %w", err)
		}
		e.Code = e.Description // Ajuste conforme necessário
		data = append(data, e)
	}
	return &data, nil
}
