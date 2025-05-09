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

// reportGestorRepository implementa o ReportRepository para a automação Gestor (ex: versão 48.2).
// Ele executa 15 queries específicas e popula models.Item com o resultado.
type reportGestorRepository struct {
	db    *sql.DB
	query querys.QueryHandler
}

// NewReportGestorRepository retorna uma instância do repositório Gestor.
func NewReportGestorRepository(db *sql.DB, q querys.QueryHandler) ReportRepository {
	return &reportGestorRepository{
		db:    db,
		query: q,
	}
}

// FindAllRangeDay executa as 15 queries para montar o relatório final (models.Item) no intervalo [startDate, endDate].
func (r *reportGestorRepository) FindAllRangeDay(startDate, endDate, store string, _ *sql.DB) (models.Item, error) {

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

	// 1) Normaliza as datas (caso precise remover HH:MM:SS e ficar só com YYYY-MM-DD)
	startSimple, endSimple, err := r.normalizeDates(startDate, endDate)
	if err != nil {
		return report, err
	}

	// (1) TOTAL_VENDIDO_DIA
	lucro, custo, valor, err := r.lucroCustoDia(startSimple, endSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar lucro/custo: %v", err)
	} else {
		report.Profit = lucro
		report.Cost = custo
		report.TotalAmount = valor
	}

	// (2) TICKET_MEDIO_DIA
	tk, err := r.ticketMedio(startSimple, endSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar ticket médio: %v", err)
	} else {
		report.AverageTicket = tk
	}

	// (3) PRODUTOS_POR_ATENDIMENTO_DIA
	ppa, err := r.produtoAtendimento(startSimple, endSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar produtos por atendimento: %v", err)
	} else {
		report.ProductsService = ppa
	}

	// (4) TOTAL_CUPONS_VALIDOS_DIA
	validCupom, err := r.totalCupomValido(startSimple, endSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar cupons válidos: %v", err)
	} else {
		report.Clients += int64(validCupom) // soma no total de clientes
	}

	// (5) TOTAL_CUPONS_CANCELADOS_DIA
	canc, err := r.cuponsCancelados(startSimple, endSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar cupons cancelados: %v", err)
	} else {
		report.CanceledCoupons = int64(canc)
	}

	// (6) VENDAS_POR_FAIXA_HORARIA
	faixa, err := r.faixaHoraria(startSimple, endSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar vendas por faixa horária: %v", err)
	} else {
		report.HourRange = *faixa
	}

	// (7) PRODUTOS_MAIS_VENDIDO_DIA
	rank, err := r.rankProdutos(startSimple, endSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar rank de produtos: %v", err)
	} else {
		report.Products = *rank
	}

	// (8) PRODUTOS_ESTATISTICA (opcional; se quiser armazenar em outro lugar)
	_, _ = r.produtosEstatistica(startSimple, endSimple, store)

	// (9) TOTAL_VENDIDO_POR_MODALIDADES
	mods, err := r.modalidadesPagamento(startSimple, endSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar modalidades de pagamento: %v", err)
	} else {
		report.Modalities = *mods
	}

	// (10) ESTORNO_DE_CUPONS_POR_OPERADOR_DIA
	estornos, err := r.estornoCupons(startSimple, endSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar estornos de cupons: %v", err)
	} else {
		report.ReversedCoupons = *estornos
	}

	// (11) DESCONTOS_DE_ITENS_DIA
	descontos, err := r.descontoItem(startSimple, endSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar descontos de itens: %v", err)
	} else {
		report.Deduction = *descontos
	}

	// (12) VENDAS_POR_SECOES_DIA
	secs, err := r.section(startSimple, endSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar vendas por seção: %v", err)
	} else {
		report.Sections = *secs
	}

	// (13) TOTAL_CUPONS_POR_OPERADOR_DIA
	cuponsOp, err := r.cuponsPorOperador(startSimple, endSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar cupons por operador: %v", err)
	} else {
		report.CouponsOperator = *cuponsOp
		for _, op := range *cuponsOp {
			report.Clients += op.Quantity
		}
	}

	// (14) VENDAS_POR_VENDEDOR_DIA
	vend, err := r.vendedores(startSimple, endSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar vendas por vendedor: %v", err)
	} else {
		report.Salesman = *vend
	}

	// (15) MODALIDADES_PAGAMENTO_OPERADOR_DIA
	modOper, err := r.modalidadeOperador(startSimple, endSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar modalidades de pagamento por operador: %v", err)
	} else {
		report.ModalitiesOperator = *modOper
	}

	return report, nil
}

//
// ░░░░░░░░░░░ FUNÇÕES AUXILIARES ░░░░░░░░░░░
//

// normalizeDates converte datas com ou sem HH:MM:SS para apenas YYYY-MM-DD
func (r *reportGestorRepository) normalizeDates(s, e string) (string, string, error) {
	layoutFull := "2006-01-02 15:04:05"
	layoutShort := "2006-01-02"

	// Tenta parsear como layoutFull
	t1, err1 := time.Parse(layoutFull, s)
	t2, err2 := time.Parse(layoutFull, e)
	if err1 == nil && err2 == nil {
		return t1.Format(layoutShort), t2.Format(layoutShort), nil
	}

	// Tenta parsear como layoutShort
	t3, err3 := time.Parse(layoutShort, s)
	t4, err4 := time.Parse(layoutShort, e)
	if err3 == nil && err4 == nil {
		return t3.Format(layoutShort), t4.Format(layoutShort), nil
	}

	return "", "", fmt.Errorf("datas '%s' e/ou '%s' em formato inválido", s, e)
}

// (1) TOTAL_VENDIDO_DIA (Lucro, Custo, Valor)
func (r *reportGestorRepository) lucroCustoDia(start, end, store string) (float64, float64, float64, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_VENDIDO_DIA())
	if err != nil {
		return 0, 0, 0, err
	}
	defer stmt.Close()

	var data struct {
		Valor float64
		Custo float64
		Lucro float64
	}
	err = stmt.QueryRow(start, end, store).Scan(&data.Valor, &data.Custo, &data.Lucro)
	if err != nil {
		return 0, 0, 0, err
	}
	return data.Lucro, data.Custo, data.Valor, nil
}

// (2) TICKET_MEDIO_DIA
func (r *reportGestorRepository) ticketMedio(start, end, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.TICKET_MEDIO_DIA())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var value float64
	err = stmt.QueryRow(store, start, end, store).Scan(&value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// (3) PRODUTOS_POR_ATENDIMENTO_DIA
func (r *reportGestorRepository) produtoAtendimento(start, end, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.PRODUTOS_POR_ATENDIMENTO_DIA())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var val float64
	err = stmt.QueryRow(store, start, end, store).Scan(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// (4) TOTAL_CUPONS_VALIDOS_DIA
func (r *reportGestorRepository) totalCupomValido(start, end, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_CUPONS_VALIDOS_DIA())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var val float64
	err = stmt.QueryRow(start, end, store).Scan(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// (5) TOTAL_CUPONS_CANCELADOS_DIA
func (r *reportGestorRepository) cuponsCancelados(start, end, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_CUPONS_CANCELADOS_DIA())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var val float64
	err = stmt.QueryRow(start, end, store).Scan(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// (6) VENDAS_POR_FAIXA_HORARIA
func (r *reportGestorRepository) faixaHoraria(start, end, store string) (*[]hour.Sales, error) {
	stmt, err := r.db.Prepare(r.query.VENDAS_POR_FAIXA_HORARIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, end, store, start, end, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []hour.Sales
	for rows.Next() {
		var e hour.Sales
		if err := rows.Scan(
			&e.Hour,
			&e.Amount,
			&e.Cupons,
			&e.Itens,
			&e.AverageItens,
			&e.AverageAmount,
		); err != nil {
			return nil, err
		}
		data = append(data, e)
	}
	return &data, nil
}

// (7) PRODUTOS_MAIS_VENDIDO_DIA
func (r *reportGestorRepository) rankProdutos(start, end, store string) (*[]product.Item, error) {
	stmt, err := r.db.Prepare(r.query.PRODUTOS_MAIS_VENDIDO_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, end, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []product.Item
	for rows.Next() {
		var e product.Item
		// Observação: a query do Gestor 48.2 retorna:
		// EAN, DESCRICAO, VALOR_UNITARIO, QUANTIDADE, TOTAL, CUSTO, TOTAL-CUSTO (Lucro)
		if err := rows.Scan(
			&e.Ean,
			&e.Description,
			&e.Price,
			&e.Quantity,
			&e.TotalAmount,
			&e.Cost,
			&e.Profit,
		); err != nil {
			return nil, err
		}
		data = append(data, e)
	}
	return &data, nil
}

// (8) PRODUTOS_ESTATISTICA (opcional)
func (r *reportGestorRepository) produtosEstatistica(start, end, store string) (*[]product.Item, error) {
	stmt, err := r.db.Prepare(r.query.PRODUTOS_ESTATISTICA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, end, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []product.Item
	for rows.Next() {
		var e product.Item
		if err := rows.Scan(
			&e.Ean,
			&e.Description,
			&e.Price,
			&e.Quantity,
			&e.TotalAmount,
			&e.Cost,
			&e.Profit,
		); err != nil {
			return nil, err
		}
		data = append(data, e)
	}
	return &data, nil
}

// (9) TOTAL_VENDIDO_POR_MODALIDADES
func (r *reportGestorRepository) modalidadesPagamento(start, end, store string) (*[]modality.Sales, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_VENDIDO_POR_MODALIDADES())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, end, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []modality.Sales
	for rows.Next() {
		var e modality.Sales
		// SELECT RD.finalizadora AS CODIGO, FZ.dsc AS DESCRICAO, COUNT(DISTINCT(RD.cupom)) AS QTD, SUM(...) AS TOTAL
		if err := rows.Scan(
			&e.Code,
			&e.Description,
			&e.Quantity,
			&e.TotalAmount,
		); err != nil {
			return nil, err
		}
		data = append(data, e)
	}
	return &data, nil
}

// (10) ESTORNO_DE_CUPONS_POR_OPERADOR_DIA
func (r *reportGestorRepository) estornoCupons(start, end, store string) (*[]reversed.Coupons, error) {
	stmt, err := r.db.Prepare(r.query.ESTORNO_DE_CUPONS_POR_OPERADOR_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, end, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []reversed.Coupons
	for rows.Next() {
		var e reversed.Coupons
		var supervisor *string
		if err := rows.Scan(
			&e.Pdv,
			&e.Date,
			&e.Hour,
			&e.Operator,
			&supervisor,
			&e.Amount,
			&e.Code,
		); err != nil {
			return nil, err
		}
		e.Authorizer = supervisor
		data = append(data, e)
	}
	return &data, nil
}

// (11) DESCONTOS_DE_ITENS_DIA
func (r *reportGestorRepository) descontoItem(start, end, store string) (*[]sales.Deduction, error) {
	stmt, err := r.db.Prepare(r.query.DESCONTOS_DE_ITENS_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, end, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []sales.Deduction
	for rows.Next() {
		var e sales.Deduction
		var autorizer *string
		if err := rows.Scan(
			&e.Date,
			&e.Hour,
			&e.Operator,
			&autorizer,
			&e.TotalDeduction,
			&e.Serial,
			&e.Item,
			&e.Ean,
			&e.Description,
		); err != nil {
			return nil, err
		}
		e.Autorizer = autorizer
		data = append(data, e)
	}
	return &data, nil
}

// (12) VENDAS_POR_SECOES_DIA
func (r *reportGestorRepository) section(start, end, store string) (*[]section.Sales, error) {
	stmt, err := r.db.Prepare(r.query.VENDAS_POR_SECOES_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, end, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []section.Sales
	for rows.Next() {
		var e section.Sales
		var qtde float64

		// Query: CODIGO, DESCRICAO, QUANTIDADE, VALOR, CUSTO, LUCRO
		if err := rows.Scan(
			&e.Code,
			&e.Description,
			&qtde,
			&e.Amout,
			&e.Cost,
			&e.Profit,
		); err != nil {
			return nil, err
		}

		// converte float -> int64
		e.Quantity = decimal.NewFromFloat(qtde).IntPart()
		data = append(data, e)
	}
	return &data, nil
}

// (13) TOTAL_CUPONS_POR_OPERADOR_DIA
func (r *reportGestorRepository) cuponsPorOperador(start, end, store string) (*[]coupon.Operator, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_CUPONS_POR_OPERADOR_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, end, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []coupon.Operator
	for rows.Next() {
		var e coupon.Operator
		var dummy float64 // se a query tiver 3 colunas e a 3a for VALOR, mas aqui só pegamos QUANTIDADE e Name
		if err := rows.Scan(&e.Quantity, &e.Name, &dummy); err != nil {
			// Se a query só retorna 2 colunas, ajuste o Scan. Depende da sua query exata.
			// Se forem só 2 colunas, comente a 'dummy'.
			// if err := rows.Scan(&e.Quantity, &e.Name); err != nil {
			return nil, err
		}
		data = append(data, e)
	}
	return &data, nil
}

// (14) VENDAS_POR_VENDEDOR_DIA
func (r *reportGestorRepository) vendedores(start, end, store string) (*[]sales.Salesman, error) {
	stmt, err := r.db.Prepare(r.query.VENDAS_POR_VENDEDOR_DIA())
	if err != nil {
		// Se a query estiver vazia, retorne slice vazio
		return &[]sales.Salesman{}, nil
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, end, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []sales.Salesman
	for rows.Next() {
		var e sales.Salesman
		// Ajuste conforme colunas da query:
		// (exemplo) DATA, NOME, VALOR, CUPONS, ITENS, MEDIA_ITENS, MEDIA_VALOR
		err := rows.Scan(
			new(string), // se quiser ignorar a DATA
			&e.Name,
			&e.Amount,
			&e.Coupons,
			&e.Itens,
			&e.AverageItems,
			&e.AverageAmount,
		)
		if err != nil {
			return nil, err
		}
		data = append(data, e)
	}
	return &data, nil
}

// (15) MODALIDADES_PAGAMENTO_OPERADOR_DIA
func (r *reportGestorRepository) modalidadeOperador(start, end, store string) (*[]modality.Operator, error) {
	stmt, err := r.db.Prepare(r.query.MODALIDADES_PAGAMENTO_OPERADOR_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, end, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []modality.Operator
	for rows.Next() {
		var e modality.Operator
		// Query: (Operador, Codigo, Descricao, Quantidade, Valor)
		if err := rows.Scan(
			&e.Operator,
			&e.Code,
			&e.Description,
			&e.Quantity,
			&e.TotalAmount,
		); err != nil {
			return nil, err
		}
		data = append(data, e)
	}
	return &data, nil
}
