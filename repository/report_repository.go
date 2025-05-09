package repository

import (
	"database/sql"
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

	"github.com/shopspring/decimal"
)

type ReportRepository interface {
	FindAllRangeDay(startDate, endDate string, store string, db *sql.DB) (models.Item, error)
}

type reportRepository struct {
	db     *sql.DB
	querys querys.QueryHandler
}

func NewReportRepository(db *sql.DB, q querys.QueryHandler) ReportRepository {
	return &reportRepository{db, q}
}

func (r *reportRepository) FindAllRangeDay(start, end, store string, db *sql.DB) (models.Item, error) {
	var report = models.Item{
		Profit:          0,
		Cost:            0,
		TotalAmount:     0,
		ProductsService: 0,
		AverageTicket:   0,
		CanceledCoupons: 0,
		Clients:         0,
	}

	lucro, custo, valor, err := r.lucroCustoDia(start, end, store)

	if err != nil {
		log.Println("Erro ao buscar o lucro e custo do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.Profit = lucro
	report.Cost = custo
	report.TotalAmount = valor

	cuponsCancelados, err := r.cuponsCancelados(start, end, store)

	if err != nil {
		log.Println("Erro ao buscar os cupons cancelados do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.CanceledCoupons = int64(cuponsCancelados)

	cuponsPorOperador, err := r.cuponsPorOperador(start, end, store)

	if err != nil {
		log.Println("Erro ao buscar os cupons por operador do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.CouponsOperator = *cuponsPorOperador

	for _, v := range *cuponsPorOperador {
		report.Clients += v.Quantity
	}

	desconto, err := r.descontoItem(start, end, store)

	if err != nil {
		log.Println("Erro ao buscar os descontos de itens do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.Deduction = *desconto

	estornoCupons, err := r.estornoCupons(start, end, store)

	if err != nil {
		log.Println("Erro ao buscar os estornos de cupons do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.ReversedCoupons = *estornoCupons

	faixaHoraria, err := r.faixaHoraria(start, end, store)

	if err != nil {
		log.Println("Erro ao buscar as vendas por faixa horaria do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.HourRange = *faixaHoraria

	modalidadeOperador, err := r.modalidadeOperador(start, end, store)

	if err != nil {
		log.Println("Erro ao buscar as modalidades de pagamento por operador do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.ModalitiesOperator = *modalidadeOperador

	modalidadesPagamento, err := r.modalidadesPagamento(start, store)

	if err != nil {
		log.Println("Erro ao buscar as modalidades de pagamento do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.Modalities = *modalidadesPagamento

	produtoAtendimento, err := r.produtoAtendimento(start, end, store)

	if err != nil {
		log.Println("Erro ao buscar os produtos por atendimento do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.ProductsService = produtoAtendimento

	rankProdutos, err := r.rankProdutos(start, store)

	if err != nil {
		log.Println("Erro ao buscar os produtos mais vendidos do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.Products = *rankProdutos

	section, err := r.section(start, end, store)

	if err != nil {
		log.Println("Erro ao buscar as vendas por seção do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.Sections = *section

	ticketMedio, err := r.ticketMedio(start, end, store)

	if err != nil {
		log.Println("Erro ao buscar o ticket medio do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.AverageTicket = ticketMedio

	vendedores, err := r.vendedores(start, end, store)

	if err != nil {
		log.Println("Erro ao buscar as vendas por vendedor do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.Salesman = *vendedores

	return report, nil
}

func (r *reportRepository) lucroCustoDia(start, end, store string) (float64, float64, float64, error) {

	var data struct {
		Lucro float64 `json:"Lucro" db:"lucro"`
		Valor float64 `json:"Valor" db:"valor"`
		Custo float64 `json:"Custo" db:"custo"`
	}

	err := r.db.QueryRow(r.querys.TOTAL_VENDIDO_DIA(), store, start, end, start, end).Scan(&data.Valor, &data.Custo, &data.Lucro)

	if err != nil {
		return 0, 0, 0, err
	}

	return data.Lucro, data.Custo, data.Valor, nil

}

func (r *reportRepository) cuponsCancelados(start, end, store string) (float64, error) {
	var value float64

	err := r.db.QueryRow(r.querys.TOTAL_CUPONS_CANCELADOS_DIA(), start, end, store).Scan(&value)

	if err != nil {
		return 0, err
	}

	return value, nil
}

func (r *reportRepository) cuponsPorOperador(start, end, store string) (*[]coupon.Operator, error) {

	rows, err := r.db.Query(r.querys.TOTAL_CUPONS_POR_OPERADOR_DIA(), start, end, store)

	if err != nil {
		return nil, err
	}

	data := []coupon.Operator{}

	for rows.Next() {

		var e coupon.Operator

		err := rows.Scan(&e.Quantity, &e.Name)

		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil
}

func (r *reportRepository) descontoItem(start, end, store string) (*[]sales.Deduction, error) {

	rows, err := r.db.Query(r.querys.DESCONTOS_DE_ITENS_DIA(), store, start, end)

	if err != nil {
		return nil, err
	}

	data := []sales.Deduction{}

	for rows.Next() {

		var e sales.Deduction

		err := rows.Scan(&e.Date, &e.Hour, &e.Operator, &e.Autorizer, &e.TotalDeduction, &e.Serial, &e.Item, &e.Ean, &e.Description)

		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil

}

func (r *reportRepository) estornoCupons(start, end, store string) (*[]reversed.Coupons, error) {

	rows, err := r.db.Query(r.querys.ESTORNO_DE_CUPONS_POR_OPERADOR_DIA(), store, start, end)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	data := []reversed.Coupons{}

	for rows.Next() {

		var e reversed.Coupons

		err := rows.Scan(&e.Pdv, &e.Date, &e.Hour, &e.Operator, &e.Authorizer, &e.Amount, &e.Code)
		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil
}

func (r *reportRepository) faixaHoraria(start, end, store string) (*[]hour.Sales, error) {

	rows, err := r.db.Query(r.querys.VENDAS_POR_FAIXA_HORARIA(), store, start, start, end)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	data := []hour.Sales{}

	for rows.Next() {

		var e hour.Sales

		err := rows.Scan(&e.Hour, &e.Amount, &e.Cupons, &e.Itens, &e.AverageItens, &e.AverageAmount)
		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil
}

func (r *reportRepository) modalidadeOperador(start, end, store string) (*[]modality.Operator, error) {

	rows, err := r.db.Query(r.querys.MODALIDADES_PAGAMENTO_OPERADOR_DIA(), store, start, end)

	if err != nil {
		return nil, err
	}

	data := []modality.Operator{}

	for rows.Next() {

		var e modality.Operator

		err := rows.Scan(&e.Operator, &e.Code, &e.Description, &e.Quantity, &e.TotalAmount)
		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil
}

func (r *reportRepository) modalidadesPagamento(start, store string) (*[]modality.Sales, error) {

	rows, err := r.db.Query(r.querys.TOTAL_VENDIDO_POR_MODALIDADES(), store, start)

	if err != nil {
		return nil, err
	}

	data := []modality.Sales{}

	for rows.Next() {

		var e modality.Sales

		err := rows.Scan(&e.Code, &e.Description, &e.Quantity, &e.TotalAmount)
		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil
}

func (r *reportRepository) produtoAtendimento(start, end, store string) (float64, error) {

	var data float64

	err := r.db.QueryRow(r.querys.PRODUTOS_POR_ATENDIMENTO_DIA(), store, start, end, start, end).Scan(&data)

	if err != nil {
		return 0, err
	}

	return data, nil
}

func (r *reportRepository) rankProdutos(start, store string) (*[]product.Item, error) {

	rows, err := r.db.Query(r.querys.PRODUTOS_MAIS_VENDIDO_DIA(), store, start)

	if err != nil {
		return nil, err
	}

	data := []product.Item{}

	for rows.Next() {

		var e product.Item

		err := rows.Scan(&e.Ean, &e.Description, &e.Quantity, &e.TotalAmount)

		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil
}

func (r *reportRepository) section(start, end, store string) (*[]section.Sales, error) {

	rows, err := r.db.Query(r.querys.VENDAS_POR_SECOES_DIA(), store, start, end)

	if err != nil {
		return nil, err
	}
	var date string

	data := []section.Sales{}

	for rows.Next() {

		var e section.Sales

		var qtde float64

		err := rows.Scan(&date, &e.Code, &e.Description, &e.Cost, &qtde, &e.Amout, &e.Profit)

		e.Quantity = decimal.NewFromFloat(qtde).IntPart()

		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil
}

func (r *reportRepository) ticketMedio(start, end, store string) (float64, error) {

	var data float64

	err := r.db.QueryRow(r.querys.TICKET_MEDIO_DIA(), store, start, end, start, end).Scan(&data)

	if err != nil {
		return 0, err
	}

	return data, nil
}

func (r *reportRepository) vendedores(start, end, store string) (*[]sales.Salesman, error) {

	rows, err := r.db.Query(r.querys.VENDAS_POR_VENDEDOR_DIA(), store, start, end, start, end)

	if err != nil {
		return nil, err
	}

	var dataRegister string

	data := []sales.Salesman{}

	for rows.Next() {

		var e sales.Salesman

		err := rows.Scan(&dataRegister, &e.Name, &e.Amount, &e.Coupons, &e.Itens, &e.AverageItems, &e.AverageAmount)

		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil
}
