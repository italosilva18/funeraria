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

// reportConcnicoRepository implementa o ReportRepository para a automação Concinco.
type reportConcnicoRepository struct {
	db    *sql.DB
	query querys.QueryHandler
}

// NewReportConcnicoRepository retorna uma instância do repositório Concinco.
func NewReportConcnicoRepository(db *sql.DB, q querys.QueryHandler) ReportRepository {
	return &reportConcnicoRepository{
		db:    db,
		query: q,
	}
}

// FindAllRangeDay executa as 15 queries para montar o relatório final (models.Item).
func (r *reportConcnicoRepository) FindAllRangeDay(startDate, endDate, store string, _ *sql.DB) (models.Item, error) {
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

	// Normaliza as datas (caso precise remover HH:MM:SS e ficar só com YYYY-MM-DD)
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
		report.Clients += int64(validCupom)
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

	// (8) PRODUTOS_ESTATISTICA (opcional)
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

// normalizeDates converte datas com ou sem HH:MM:SS para apenas YYYY-MM-DD
func (r *reportConcnicoRepository) normalizeDates(s, e string) (string, string, error) {
	layoutFull := "2006-01-02 15:04:05"
	layoutShort := "2006-01-02"

	// Tenta parsear como layoutFull
	t1, err1 := time.Parse(layoutFull, s)
	t2, err2 := time.Parse(layoutFull, e)
	if err1 == nil && err2 == nil {
		// Converte para formato Oracle DD-MM-YYYY
		start := t1.Format("02-01-2006")
		end := t2.Format("02-01-2006")
		return start, end, nil
	}

	// Tenta parsear como layoutShort
	t3, err3 := time.Parse(layoutShort, s)
	t4, err4 := time.Parse(layoutShort, e)
	if err3 == nil && err4 == nil {
		// Converte para formato Oracle DD-MM-YYYY
		start := t3.Format("02-01-2006")
		end := t4.Format("02-01-2006")
		return start, end, nil
	}

	return "", "", fmt.Errorf("datas '%s' e/ou '%s' em formato inválido", s, e)
}

// (1) TOTAL_VENDIDO_DIA - params: start (DD-MM-YYYY), store
func (r *reportConcnicoRepository) lucroCustoDia(start, end, store string) (float64, float64, float64, error) {
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
	// Ordem: :1 = data, :2 = empresa
	err = stmt.QueryRow(start, store).Scan(&data.Valor, &data.Custo, &data.Lucro)
	if err != nil {
		return 0, 0, 0, err
	}
	return data.Lucro, data.Custo, data.Valor, nil
}

// (2) TICKET_MEDIO_DIA - params: store, start (DD-MM-YYYY)
func (r *reportConcnicoRepository) ticketMedio(start, end, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.TICKET_MEDIO_DIA())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var value float64
	// Ordem: :1 = empresa, :2 = data
	err = stmt.QueryRow(store, start).Scan(&value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// (3) PRODUTOS_POR_ATENDIMENTO_DIA - params: store, start (DD-MM-YYYY)
func (r *reportConcnicoRepository) produtoAtendimento(start, end, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.PRODUTOS_POR_ATENDIMENTO_DIA())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var val float64
	// Ordem: :1 = empresa, :2 = data
	err = stmt.QueryRow(store, start).Scan(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// (4) TOTAL_CUPONS_VALIDOS_DIA - params: store, start (DD-MM-YYYY)
func (r *reportConcnicoRepository) totalCupomValido(start, end, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_CUPONS_VALIDOS_DIA())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var val float64
	// Ordem: :1 = empresa, :2 = data
	err = stmt.QueryRow(store, start).Scan(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// (5) TOTAL_CUPONS_CANCELADOS_DIA - params: start (DD-MM-YYYY), store
func (r *reportConcnicoRepository) cuponsCancelados(start, end, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_CUPONS_CANCELADOS_DIA())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var val float64
	// Ordem: :1 = data, :2 = empresa
	err = stmt.QueryRow(start, store).Scan(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// (6) VENDAS_POR_FAIXA_HORARIA - params: store, start (DD-MM-YYYY)
func (r *reportConcnicoRepository) faixaHoraria(start, end, store string) (*[]hour.Sales, error) {
	stmt, err := r.db.Prepare(r.query.VENDAS_POR_FAIXA_HORARIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Ordem: :1 = empresa, :2 = data
	rows, err := stmt.Query(store, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []hour.Sales
	for rows.Next() {
		var e hour.Sales
		// Scan: HORA, VALOR, NUM_CUPONS, NUM_ITENS, MEDIA_ITENS, MEDIA_VALOR
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

// (7) PRODUTOS_MAIS_VENDIDO_DIA - params: store, start (DD-MM-YYYY)
func (r *reportConcnicoRepository) rankProdutos(start, end, store string) (*[]product.Item, error) {
	stmt, err := r.db.Prepare(r.query.PRODUTOS_MAIS_VENDIDO_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Ordem: :1 = empresa, :2 = data
	rows, err := stmt.Query(store, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []product.Item
	for rows.Next() {
		var e product.Item
		// Scan: EAN, DESCRICAO, VALOR_UNITARIO, QUANTIDADE, TOTAL, CUSTO, LUCRO
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

// (8) PRODUTOS_ESTATISTICA - params: store, start (DD-MM-YYYY)
func (r *reportConcnicoRepository) produtosEstatistica(start, end, store string) (*[]product.Item, error) {
	stmt, err := r.db.Prepare(r.query.PRODUTOS_ESTATISTICA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Ordem: :1 = empresa, :2 = data
	rows, err := stmt.Query(store, start)
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

// (9) TOTAL_VENDIDO_POR_MODALIDADES - params: store, start (DD-MM-YYYY)
func (r *reportConcnicoRepository) modalidadesPagamento(start, end, store string) (*[]modality.Sales, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_VENDIDO_POR_MODALIDADES())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Ordem: :1 = empresa, :2 = data
	rows, err := stmt.Query(store, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []modality.Sales
	for rows.Next() {
		var e modality.Sales
		// Scan: CODIGO, DESCRICAO, QUANTIDADE, TOTAL
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

// (10) ESTORNO_DE_CUPONS_POR_OPERADOR_DIA - params: start (DD-MM-YYYY), store
func (r *reportConcnicoRepository) estornoCupons(start, end, store string) (*[]reversed.Coupons, error) {
	stmt, err := r.db.Prepare(r.query.ESTORNO_DE_CUPONS_POR_OPERADOR_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Ordem: :1 = data, :2 = empresa
	rows, err := stmt.Query(start, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []reversed.Coupons
	for rows.Next() {
		var e reversed.Coupons
		var supervisor *string
		// Scan: PDV, DATA, HORA, OPERADOR, SUPERVISOR, VALOR, CUPOM
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

// (11) DESCONTOS_DE_ITENS_DIA - params: start (DD-MM-YYYY), store
func (r *reportConcnicoRepository) descontoItem(start, end, store string) (*[]sales.Deduction, error) {
	stmt, err := r.db.Prepare(r.query.DESCONTOS_DE_ITENS_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Ordem: :1 = data, :2 = empresa
	rows, err := stmt.Query(start, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []sales.Deduction
	for rows.Next() {
		var e sales.Deduction
		var autorizer *string
		// Scan: DATA, HORA, OPERADOR, SUPERVISOR, VALOR, CUPOM, ITEM, EAN, DESCRICAO
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

// (12) VENDAS_POR_SECOES_DIA - params: start (DD-MM-YYYY), store
func (r *reportConcnicoRepository) section(start, end, store string) (*[]section.Sales, error) {
	stmt, err := r.db.Prepare(r.query.VENDAS_POR_SECOES_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Ordem: :1 = data, :2 = empresa
	rows, err := stmt.Query(start, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []section.Sales
	for rows.Next() {
		var e section.Sales
		var qtde float64
		// Scan: CODIGO, DESCRICAO, QUANTIDADE, VALOR, CUSTO, LUCRO
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
		// Converte float -> int64
		e.Quantity = decimal.NewFromFloat(qtde).IntPart()
		data = append(data, e)
	}
	return &data, nil
}

// (13) TOTAL_CUPONS_POR_OPERADOR_DIA - params: store, start (DD-MM-YYYY)
func (r *reportConcnicoRepository) cuponsPorOperador(start, end, store string) (*[]coupon.Operator, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_CUPONS_POR_OPERADOR_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Ordem: :1 = empresa, :2 = data
	rows, err := stmt.Query(store, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []coupon.Operator
	for rows.Next() {
		var e coupon.Operator
		var dummy float64 // VALOR é retornado mas não usado no modelo
		// Scan: QTD_CUPONS, OPERADOR, VALOR
		if err := rows.Scan(&e.Quantity, &e.Name, &dummy); err != nil {
			return nil, err
		}
		data = append(data, e)
	}
	return &data, nil
}

// (14) VENDAS_POR_VENDEDOR_DIA - params: start (DD-MM-YYYY), store
func (r *reportConcnicoRepository) vendedores(start, end, store string) (*[]sales.Salesman, error) {
	stmt, err := r.db.Prepare(r.query.VENDAS_POR_VENDEDOR_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Ordem: :1 = data, :2 = empresa
	rows, err := stmt.Query(start, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []sales.Salesman
	for rows.Next() {
		var e sales.Salesman
		var dataStr string    // DATA é retornada mas não usada
		var itemFloat float64 // ITEM pode ser decimal
		// Scan: DATA, VENDEDOR, VALOR, CUPONS, ITEM, MEDIA_ITENS, MEDIA_VALOR
		if err := rows.Scan(
			&dataStr,
			&e.Name,
			&e.Amount,
			&e.Coupons,
			&itemFloat, // Mudança aqui
			&e.AverageItems,
			&e.AverageAmount,
		); err != nil {
			return nil, err
		}
		// Converte float para int64
		e.Itens = int64(itemFloat)
		data = append(data, e)
	}
	return &data, nil
}

// (15) MODALIDADES_PAGAMENTO_OPERADOR_DIA - params: store, start (DD-MM-YYYY)
func (r *reportConcnicoRepository) modalidadeOperador(start, end, store string) (*[]modality.Operator, error) {
	stmt, err := r.db.Prepare(r.query.MODALIDADES_PAGAMENTO_OPERADOR_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Ordem: :1 = empresa, :2 = data
	rows, err := stmt.Query(store, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []modality.Operator
	for rows.Next() {
		var e modality.Operator
		var dataStr string // DATA é retornada mas não usada
		// Scan: DATA, OPERADOR, COD_PAGTO, DESCRICAO, QTD, VALOR
		if err := rows.Scan(
			&dataStr,
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
