package repository

import (
	"database/sql"
	"fmt"
	"log"
	"margem/robo/helpers"
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

// TODO Refatorar e colocar um design para a interface
type reportRmsRepository struct {
	db     *sql.DB
	querys querys.QueryHandler
}

func NewreportRmsRepository(db *sql.DB, q querys.QueryHandler) ReportRepository {
	return &reportRmsRepository{db, q}
}

func (r *reportRmsRepository) FindAllRangeDay(start, end, store string, db *sql.DB) (models.Item, error) {

	dateFormatSimple, _ := helpers.DateTo1AAMMDD(start)

	var report = models.Item{
		Profit:          0,
		Cost:            0,
		TotalAmount:     0,
		ProductsService: 0,
		AverageTicket:   0,
		CanceledCoupons: 0,
		Clients:         0,
	}

	lucro, custo, valor, err := r.lucroCustoDia(dateFormatSimple, store)

	if err != nil {
		log.Println("Erro ao buscar o lucro e custo do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.Profit = lucro
	report.Cost = custo
	report.TotalAmount = valor

	cuponsCancelados, err := r.cuponsCancelados(dateFormatSimple, store)

	if err != nil {
		log.Println("Erro ao buscar os cupons cancelados do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.CanceledCoupons = int64(cuponsCancelados)

	cuponsPorOperador, err := r.cuponsPorOperador(dateFormatSimple, store)

	if err != nil {
		log.Println("Erro ao buscar os cupons por operador do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.CouponsOperator = *cuponsPorOperador

	for _, v := range *cuponsPorOperador {
		report.Clients += v.Quantity
	}

	// TODO Refatorar para pegar o desconto de itens depois da query está resolvida
	desconto, err := r.descontoItem(dateFormatSimple, store)

	if err != nil {
		log.Println("Erro ao buscar os descontos de itens do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.Deduction = *desconto

	// TODO Refatorar para pegar os estornos de cupons depois da query está resolvida
	estornoCupons, err := r.estornoCupons(dateFormatSimple, store)

	if err != nil {
		log.Println("Erro ao buscar os estornos de cupons do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.ReversedCoupons = *estornoCupons

	faixaHoraria, err := r.faixaHoraria(dateFormatSimple, store)

	if err != nil {
		log.Println("Erro ao buscar as vendas por faixa horaria do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.HourRange = *faixaHoraria

	modalidadeOperador, err := r.modalidadeOperador(dateFormatSimple, store)

	if err != nil {
		log.Println("Erro ao buscar as modalidades de pagamento por operador do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.ModalitiesOperator = *modalidadeOperador

	modalidadesPagamento, err := r.modalidadesPagamento(dateFormatSimple, store)

	if err != nil {
		log.Println("Erro ao buscar as modalidades de pagamento do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.Modalities = *modalidadesPagamento

	produtoAtendimento, err := r.produtoAtendimento(dateFormatSimple, store)

	if err != nil {
		log.Println("Erro ao buscar os produtos por atendimento do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.ProductsService = produtoAtendimento

	rankProdutos, err := r.rankProdutos(dateFormatSimple, store)

	if err != nil {
		log.Println("Erro ao buscar os produtos mais vendidos do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.Products = *rankProdutos

	section, err := r.section(dateFormatSimple, store)

	if err != nil {
		log.Println("Erro ao buscar as vendas por seção do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.Sections = *section

	ticketMedio, err := r.ticketMedio(dateFormatSimple, store)

	if err != nil {
		log.Println("Erro ao buscar o ticket medio do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.AverageTicket = ticketMedio

	vendedores, err := r.vendedores(dateFormatSimple, store)

	if err != nil {
		log.Println("Erro ao buscar as vendas por vendedor do dia")
		log.Printf("Erro: %v \n", err)
	}

	report.Salesman = *vendedores

	return report, nil
}

// DONE query concluída
func (r *reportRmsRepository) lucroCustoDia(start, store string) (float64, float64, float64, error) {

	var data struct {
		Lucro float64 `json:"Lucro" db:"lucro"`
		Valor float64 `json:"Valor" db:"valor"`
		Custo float64 `json:"Custo" db:"custo"`
	}

	stm, _ := r.db.Prepare(r.querys.TOTAL_VENDIDO_DIA())

	err := stm.QueryRow(start, store).Scan(&data.Valor, &data.Custo, &data.Lucro)

	if err != nil {
		return 0, 0, 0, err
	}

	return data.Lucro, data.Custo, data.Valor, nil
}

// DONE query concluída
func (r *reportRmsRepository) cuponsCancelados(start, store string) (float64, error) {
	var value float64

	stm, err := r.db.Prepare(r.querys.TOTAL_CUPONS_CANCELADOS_DIA())

	if err != nil {
		return 0, err
	}

	err = stm.QueryRow(start, store).Scan(&value)

	if err != nil {
		return 0, err
	}

	return value, nil
}

// DONE query concluída
func (r *reportRmsRepository) cuponsPorOperador(start, store string) (*[]coupon.Operator, error) {

	stm, err := r.db.Prepare(r.querys.TOTAL_CUPONS_POR_OPERADOR_DIA())

	if err != nil {
		return nil, err
	}

	rows, err := stm.Query(start, store)

	if err != nil {
		return nil, err
	}

	data := []coupon.Operator{}

	for rows.Next() {

		var e coupon.Operator
		var value float64

		err := rows.Scan(&e.Quantity, &e.Name, &value)

		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil
}

// DONE query concluída
func (r *reportRmsRepository) descontoItem(start, store string) (*[]sales.Deduction, error) {

	stm, err := r.db.Prepare(r.querys.DESCONTOS_DE_ITENS_DIA())

	if err != nil {
		return nil, err
	}

	rows, err := stm.Query(start, store)

	if err != nil {
		return nil, err
	}

	data := []sales.Deduction{}
	var nullField *string

	for rows.Next() {

		var e sales.Deduction

		err := rows.Scan(&e.Date, &e.Hour, &e.Operator, &nullField, &e.TotalDeduction, &e.Serial, &e.Item, &e.Ean, &e.Description)

		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil
}

// DONE query com erro
func (r *reportRmsRepository) estornoCupons(start, store string) (*[]reversed.Coupons, error) {

	stm, err := r.db.Prepare(r.querys.ESTORNO_DE_CUPONS_POR_OPERADOR_DIA())

	if err != nil {
		return nil, err
	}

	rows, err := stm.Query(start, store)

	var item *string

	if err != nil {
		log.Println(err)
		return nil, err
	}

	data := []reversed.Coupons{}

	for rows.Next() {

		var e reversed.Coupons

		err := rows.Scan(&e.Pdv, &e.Date, &e.Hour, &e.Operator, &item, &e.Amount, &e.Code)
		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil
}

// DONE query concluída
func (r *reportRmsRepository) faixaHoraria(start, store string) (*[]hour.Sales, error) {

	stm, err := r.db.Prepare(r.querys.VENDAS_POR_FAIXA_HORARIA())

	if err != nil {
		return nil, err
	}

	rows, err := stm.Query(start, store)

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

// DONE query concluída
func (r *reportRmsRepository) modalidadeOperador(start, store string) (*[]modality.Operator, error) {

	stm, err := r.db.Prepare(r.querys.MODALIDADES_PAGAMENTO_OPERADOR_DIA())

	if err != nil {
		return nil, err
	}

	rows, err := stm.Query(store, start)

	if err != nil {
		return nil, err
	}

	data := []modality.Operator{}

	for rows.Next() {

		var e modality.Operator
		var dataString string

		err := rows.Scan(&dataString, &e.Operator, &e.Code, &e.Description, &e.Quantity, &e.TotalAmount)
		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil
}

// DONE query concluída
func (r *reportRmsRepository) modalidadesPagamento(start, store string) (*[]modality.Sales, error) {

	stm, err := r.db.Prepare(r.querys.TOTAL_VENDIDO_POR_MODALIDADES())

	if err != nil {
		return nil, err
	}

	rows, err := stm.Query(store, start)

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

// DONE query concluída
func (r *reportRmsRepository) produtoAtendimento(start, store string) (float64, error) {
	var data float64

	stm, err := r.db.Prepare(r.querys.PRODUTOS_POR_ATENDIMENTO_DIA())
	if err != nil {
		return 0, err
	}
	defer stm.Close()

	err = stm.QueryRow(start, store).Scan(&data)
	if err != nil {
		return 0, err
	}

	return data, nil
}

// DONE query concluída
func (r *reportRmsRepository) rankProdutos(start, store string) (*[]product.Item, error) {

	stm, err := r.db.Prepare(r.querys.PRODUTOS_MAIS_VENDIDO_DIA())

	if err != nil {
		return nil, err
	}

	rows, err := stm.Query(start, store)

	if err != nil {
		return nil, err
	}

	data := []product.Item{}

	for rows.Next() {

		var e product.Item

		err := rows.Scan(&e.Ean, &e.Description, &e.Price, &e.Quantity, &e.TotalAmount, &e.Cost, &e.Profit)

		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil
}

// DONE query concluída
func (r *reportRmsRepository) section(start, store string) (*[]section.Sales, error) {

	stm, err := r.db.Prepare(r.querys.VENDAS_POR_SECOES_DIA())

	if err != nil {
		return nil, err
	}

	rows, err := stm.Query(start, store)

	if err != nil {
		return nil, err
	}
	//var date string

	data := []section.Sales{}

	for rows.Next() {

		var e section.Sales

		var qtde float64

		err := rows.Scan(&e.Code, &e.Description, &qtde, &e.Amout)

		e.Quantity = decimal.NewFromFloat(qtde).IntPart()

		if err != nil {
			return nil, err
		}

		data = append(data, e)
	}

	return &data, nil
}

// DONE query concluída
func (r *reportRmsRepository) ticketMedio(start, store string) (float64, error) {
	var data float64

	// Preparar a consulta SQL
	stm, err := r.db.Prepare(r.querys.TICKET_MEDIO_DIA())
	if err != nil {
		return 0, fmt.Errorf("error preparing statement: %v", err)
	}
	defer stm.Close() // Garantir que o statement será fechado após o uso

	// Log para depuração
	//log.Printf("Executing query with start: %s, store: %s", start, store)

	// Executar a consulta e escanear o resultado
	err = stm.QueryRow(start, store).Scan(&data)
	if err != nil {
		return 0, fmt.Errorf("error executing query: %v", err)
	}

	return data, nil
}

// DONE query faltando
func (r *reportRmsRepository) vendedores(start, store string) (*[]sales.Salesman, error) {

	stm, err := r.db.Prepare(r.querys.VENDAS_POR_VENDEDOR_DIA())

	if err != nil {
		return nil, err
	}

	var dataRegister int64
	var cupoms float64

	data := []sales.Salesman{}

	rows, err := stm.Query(start, store)

	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var e sales.Salesman

		err := rows.Scan(&dataRegister, &e.Name, &cupoms, &e.Coupons, &e.Itens, &e.AverageItems, &e.AverageAmount)

		if err != nil {
			return nil, err
		}

		e.Amount = float64(e.Itens)

		data = append(data, e)
	}

	return &data, nil
}
