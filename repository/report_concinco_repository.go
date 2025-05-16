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

	// Normaliza a data (no Concinco usamos apenas startDate, endDate é ignorado)
	startSimple, err := r.normalizeDate(startDate)
	if err != nil {
		return report, err
	}

	// (1) TOTAL_VENDIDO_DIA - params: date (TO_DATE), store
	lucro, custo, valor, err := r.lucroCustoDia(startSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar lucro/custo: %v", err)
	} else {
		report.Profit = lucro
		report.Cost = custo
		report.TotalAmount = valor
	}

	// (2) TICKET_MEDIO_DIA - params: date, store
	tk, err := r.ticketMedio(startSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar ticket médio: %v", err)
	} else {
		report.AverageTicket = tk
	}

	// (3) PRODUTOS_POR_ATENDIMENTO_DIA - params: date, store
	ppa, err := r.produtoAtendimento(startSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar produtos por atendimento: %v", err)
	} else {
		report.ProductsService = ppa
	}

	// (4) TOTAL_CUPONS_VALIDOS_DIA - params: date, store
	validCupom, err := r.totalCupomValido(startSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar cupons válidos: %v", err)
	} else {
		report.Clients += int64(validCupom)
	}

	// (5) TOTAL_CUPONS_CANCELADOS_DIA - params: date, store
	canc, err := r.cuponsCancelados(startSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar cupons cancelados: %v", err)
	} else {
		report.CanceledCoupons = int64(canc)
	}

	// (6) VENDAS_POR_FAIXA_HORARIA - params: date, store
	faixa, err := r.faixaHoraria(startSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar vendas por faixa horária: %v", err)
	} else {
		report.HourRange = *faixa
	}

	// (7) PRODUTOS_MAIS_VENDIDO_DIA - params: date, store
	rank, err := r.rankProdutos(startSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar rank de produtos: %v", err)
	} else {
		report.Products = *rank
	}

	// (8) PRODUTOS_ESTATISTICA - params: date, store
	_, _ = r.produtosEstatistica(startSimple, store)

	// (9) TOTAL_VENDIDO_POR_MODALIDADES - params: date, store
	mods, err := r.modalidadesPagamento(startSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar modalidades de pagamento: %v", err)
	} else {
		report.Modalities = *mods
	}

	// (10) ESTORNO_DE_CUPONS_POR_OPERADOR_DIA - params: date, store
	estornos, err := r.estornoCupons(startSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar estornos de cupons: %v", err)
	} else {
		report.ReversedCoupons = *estornos
	}

	// (11) DESCONTOS_DE_ITENS_DIA - params: date, store
	descontos, err := r.descontoItem(startSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar descontos de itens: %v", err)
	} else {
		report.Deduction = *descontos
	}

	// (12) VENDAS_POR_SECOES_DIA - params: date, store (nota: query usa :1 duas vezes)
	secs, err := r.section(startSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar vendas por seção: %v", err)
	} else {
		report.Sections = *secs
	}

	// (13) TOTAL_CUPONS_POR_OPERADOR_DIA - params: date, store
	cuponsOp, err := r.cuponsPorOperador(startSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar cupons por operador: %v", err)
	} else {
		report.CouponsOperator = *cuponsOp
		for _, op := range *cuponsOp {
			report.Clients += op.Quantity
		}
	}

	// (14) VENDAS_POR_VENDEDOR_DIA - params: date, store
	vend, err := r.vendedores(startSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar vendas por vendedor: %v", err)
	} else {
		report.Salesman = *vend
	}

	// (15) MODALIDADES_PAGAMENTO_OPERADOR_DIA - params: date, store
	modOper, err := r.modalidadeOperador(startSimple, store)
	if err != nil {
		log.Printf("Erro ao buscar modalidades de pagamento por operador: %v", err)
	} else {
		report.ModalitiesOperator = *modOper
	}

	return report, nil
}

// normalizeDate converte data com ou sem HH:MM:SS para formato Oracle DD/MM/YYYY
func (r *reportConcnicoRepository) normalizeDate(s string) (string, error) {
	layoutFull := "2006-01-02 15:04:05"
	layoutShort := "2006-01-02"

	// Tenta parsear como layoutFull
	t1, err1 := time.Parse(layoutFull, s)
	if err1 == nil {
		return t1.Format("02/01/2006"), nil
	}

	// Tenta parsear como layoutShort
	t2, err2 := time.Parse(layoutShort, s)
	if err2 == nil {
		return t2.Format("02/01/2006"), nil
	}

	return "", fmt.Errorf("data '%s' em formato inválido", s)
}

// (1) TOTAL_VENDIDO_DIA - ORDER: date, store
func (r *reportConcnicoRepository) lucroCustoDia(start, store string) (float64, float64, float64, error) {
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
	// Ordem: :1 = date (TO_DATE(:1, 'DD/MM/YYYY')), :2 = store
	err = stmt.QueryRow(start, store).Scan(&data.Valor, &data.Custo, &data.Lucro)
	if err != nil {
		return 0, 0, 0, err
	}
	return data.Lucro, data.Custo, data.Valor, nil
}

// (2) TICKET_MEDIO_DIA - ORDER: date, store
func (r *reportConcnicoRepository) ticketMedio(start, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.TICKET_MEDIO_DIA())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var value float64
	err = stmt.QueryRow(start, store).Scan(&value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// (3) PRODUTOS_POR_ATENDIMENTO_DIA - ORDER: date, store
func (r *reportConcnicoRepository) produtoAtendimento(start, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.PRODUTOS_POR_ATENDIMENTO_DIA())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var val float64
	err = stmt.QueryRow(start, store).Scan(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// (4) TOTAL_CUPONS_VALIDOS_DIA - ORDER: date, store
func (r *reportConcnicoRepository) totalCupomValido(start, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_CUPONS_VALIDOS_DIA())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var val float64
	err = stmt.QueryRow(start, store).Scan(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// (5) TOTAL_CUPONS_CANCELADOS_DIA - ORDER: date, store
func (r *reportConcnicoRepository) cuponsCancelados(start, store string) (float64, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_CUPONS_CANCELADOS_DIA())
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var val float64
	err = stmt.QueryRow(start, store).Scan(&val)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// (6) VENDAS_POR_FAIXA_HORARIA - ORDER: date, store
func (r *reportConcnicoRepository) faixaHoraria(start, store string) (*[]hour.Sales, error) {
	stmt, err := r.db.Prepare(r.query.VENDAS_POR_FAIXA_HORARIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, store)
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

// (7) PRODUTOS_MAIS_VENDIDO_DIA - ORDER: date, store
func (r *reportConcnicoRepository) rankProdutos(start, store string) (*[]product.Item, error) {
	stmt, err := r.db.Prepare(r.query.PRODUTOS_MAIS_VENDIDO_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, store)
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

// (8) PRODUTOS_ESTATISTICA - ORDER: date, store
func (r *reportConcnicoRepository) produtosEstatistica(start, store string) (*[]product.Item, error) {
	stmt, err := r.db.Prepare(r.query.PRODUTOS_ESTATISTICA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, store)
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

// (9) TOTAL_VENDIDO_POR_MODALIDADES - ORDER: date, store
func (r *reportConcnicoRepository) modalidadesPagamento(start, store string) (*[]modality.Sales, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_VENDIDO_POR_MODALIDADES())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, store)
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

// (10) ESTORNO_DE_CUPONS_POR_OPERADOR_DIA - ORDER: date, store
func (r *reportConcnicoRepository) estornoCupons(start, store string) (*[]reversed.Coupons, error) {
	stmt, err := r.db.Prepare(r.query.ESTORNO_DE_CUPONS_POR_OPERADOR_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

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

// (11) DESCONTOS_DE_ITENS_DIA - ORDER: date, store
func (r *reportConcnicoRepository) descontoItem(start, store string) (*[]sales.Deduction, error) {
	stmt, err := r.db.Prepare(r.query.DESCONTOS_DE_ITENS_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, store)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []sales.Deduction
	for rows.Next() {
		var e sales.Deduction
		var supervisor *string
		// Scan: DATA, HORA, OPERADOR, SUPERVISOR, VALOR, CUPOM, ITEM, EAN, DESCRICAO
		if err := rows.Scan(
			&e.Date,
			&e.Hour,
			&e.Operator,
			&supervisor,
			&e.TotalDeduction,
			&e.Serial,
			&e.Item,
			&e.Ean,
			&e.Description,
		); err != nil {
			return nil, err
		}
		e.Autorizer = supervisor
		data = append(data, e)
	}
	return &data, nil
}

// (12) VENDAS_POR_SECOES_DIA - ORDER: date, store (nota: query usa :1 duas vezes)
func (r *reportConcnicoRepository) section(start, store string) (*[]section.Sales, error) {
	stmt, err := r.db.Prepare(r.query.VENDAS_POR_SECOES_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// A query usa :1 duas vezes (na CTE e na query principal)
	rows, err := stmt.Query(start, store, start, store)
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

// (13) TOTAL_CUPONS_POR_OPERADOR_DIA - ORDER: date, store
func (r *reportConcnicoRepository) cuponsPorOperador(start, store string) (*[]coupon.Operator, error) {
	stmt, err := r.db.Prepare(r.query.TOTAL_CUPONS_POR_OPERADOR_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, store)
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

// (14) VENDAS_POR_VENDEDOR_DIA - ORDER: date, store
func (r *reportConcnicoRepository) vendedores(start, store string) (*[]sales.Salesman, error) {
	stmt, err := r.db.Prepare(r.query.VENDAS_POR_VENDEDOR_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

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
			&itemFloat,
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

// (15) MODALIDADES_PAGAMENTO_OPERADOR_DIA - ORDER: date, store
func (r *reportConcnicoRepository) modalidadeOperador(start, store string) (*[]modality.Operator, error) {
	stmt, err := r.db.Prepare(r.query.MODALIDADES_PAGAMENTO_OPERADOR_DIA())
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(start, store)
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
