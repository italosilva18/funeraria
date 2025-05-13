package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"margem/robo/connections"
	"margem/robo/models"
	"margem/robo/models/config"
	"margem/robo/querys"
	"margem/robo/querys/arius"
	"margem/robo/querys/arius_erp"
	"margem/robo/querys/concinco"
	"margem/robo/querys/derevo"
	"margem/robo/querys/gestor"
	"margem/robo/querys/rms"
	"margem/robo/repository"
	"net/http"
	"strconv"
	"time"
	_ "time/tzdata"
)

// ParseService expõe o método Parse para sincronizar dados.
type ParseService interface {
	Parse(startDate, endDate string) error
}

// parseService implementa ParseService e armazena informações da loja e o modo de log.
type parseService struct {
	loja    config.Loja
	modoLog string // "detalhado" ou "resumido"
}

// NewParseService retorna uma instância de parseService.
func NewParseService(loja config.Loja, modoLog string) ParseService {
	return &parseService{
		loja:    loja,
		modoLog: modoLog,
	}
}

// Parse executa o fluxo completo de sincronização:
// 1) Estabelece a conexão com o banco
// 2) Instancia o repositório conforme a automação/versão
// 3) Executa a consulta para gerar o relatório
// 4) Preenche os dados fixos do relatório com informações da loja
// 5) Envia o relatório para a API externa
func (p *parseService) Parse(startDate, endDate string) error {
	p.logDetalhado("🔎 [Parse] Iniciando coleta de dados para a loja '%s' (Loja #%d) do período '%s' até '%s'...", p.loja.NomeFantasia, p.loja.NumeroLoja, startDate, endDate)

	// 1. Estabelece a conexão com o banco de dados
	dbHandler, err := connections.GetConnection(p.loja.Automacao, p.loja)
	if err != nil {
		return fmt.Errorf("erro ao obter driver de banco: %w", err)
	}
	db, err := dbHandler.DB()
	if err != nil {
		return fmt.Errorf("erro ao abrir conexão com o banco: %w", err)
	}
	defer db.Close()
	p.logDetalhado("💾 [Parse] Conexão BD estabelecida com sucesso para a loja '%s'.", p.loja.NomeFantasia)

	// 2. Obtenção do QueryHandler e criação do repositório conforme a automação.
	var queryHandler interface{}
	var repo repository.ReportRepository

	switch p.loja.Automacao {
	case "RMS":
		queryHandler = rms.GetQueryHandler(fmt.Sprintf("%.1f", p.loja.CodigoSistema))
		repo = repository.NewreportRmsRepository(db, queryHandler.(querys.QueryHandler))
	case "GESTOR":
		queryHandler = gestor.GetQueryHandler(fmt.Sprintf("%.1f", p.loja.CodigoSistema))
		repo = repository.NewReportGestorRepository(db, queryHandler.(querys.QueryHandler))
	case "ARIUS":
		queryHandler = arius.GetQueryHandler(fmt.Sprintf("%.1f", p.loja.CodigoSistema))
		repo = repository.NewReportGestorRepository(db, queryHandler.(querys.QueryHandler))
	case "ARIUS-ERP":
		queryHandler = arius_erp.GetQueryHandler(fmt.Sprintf("%.1f", p.loja.CodigoSistema))
		repo = repository.NewReportGestorRepository(db, queryHandler.(querys.QueryHandler))
	case "DEREVO":
		queryHandler = derevo.GetQueryHandler(fmt.Sprintf("%.1f", p.loja.CodigoSistema))
		repo = repository.NewReportDerevoRepository(db, queryHandler.(querys.QueryHandler))
	case "CONCINCO":
		queryHandler = concinco.GetQueryHandler(fmt.Sprintf("%.1f", p.loja.CodigoSistema))
		repo = repository.NewReportConcnicoRepository(db, queryHandler.(querys.QueryHandler))
	default:
		return fmt.Errorf("automação '%s' não suportada", p.loja.Automacao)
	}
	p.logDetalhado("⚙️ [Parse] Handler de queries para '%s' (Sistema %.1f) obtido com sucesso.", p.loja.Automacao, p.loja.CodigoSistema)

	// 3. Executa a consulta para gerar o relatório
	report, err := repo.FindAllRangeDay(startDate, endDate, strconv.Itoa(p.loja.NumeroLoja), db)
	if err != nil {
		log.Printf("❌ [Parse] Erro ao buscar dados da loja '%s': %v", p.loja.NomeFantasia, err)
		return err
	}
	p.logDetalhado("✅ [Parse] Consulta de dados concluída para a loja '%s'.", p.loja.NomeFantasia)

	// 4. Preenche os dados fixos do relatório com informações da loja
	p.fillStaticData(&report)

	// Ajusta datas, timezone e versão do relatório
	loc, errLoc := time.LoadLocation("America/Sao_Paulo")
	if errLoc != nil {
		loc = time.UTC
	}
	dataReferencia, errTime := time.ParseInLocation("2006-01-02 15:04:05", startDate, loc)
	if errTime != nil {
		p.logDetalhado("[Parse] Aviso: Não foi possível parsear a data '%s' com timezone. Prosseguindo...", startDate)
	}
	report.Date = dataReferencia.Format("02/01/2006")
	report.CreateAt = time.Date(dataReferencia.Year(), dataReferencia.Month(), dataReferencia.Day(), 0, 0, 0, 0, loc)
	report.Synce = time.Now().In(loc)
	report.Delete = time.Now().In(loc).AddDate(1, 0, 0)
	report.Version = "11.0.0" // Versão arbitrária para envio

	// 5. Envia o relatório para a API externa
	if err := p.sendData(report); err != nil {
		log.Printf("❌ [Parse] Erro ao enviar dados da loja '%s': %v", p.loja.NomeFantasia, err)
		return err
	}
	p.logDetalhado("📡 [Parse] Dados enviados para a API com sucesso para a loja '%s'.", p.loja.NomeFantasia)

	// Log de conclusão
	if p.modoLog == "resumido" {
		log.Printf("✅ [Parse] Loja '%s' (#%d) -> data %s sincronizada com sucesso.", p.loja.NomeFantasia, p.loja.NumeroLoja, startDate[:10])
	} else {
		p.logDetalhado("✅ [Parse] Finalizado com sucesso para a loja '%s'.", p.loja.NomeFantasia)
	}

	return nil
}

// fillStaticData preenche campos fixos do relatório com base na configuração da loja.
func (p *parseService) fillStaticData(report *models.Item) {
	report.Store = p.loja.NomeFantasia
	report.Serial = p.loja.Licenca
	report.CNPJ = p.loja.Cnpj
	report.State = p.loja.Localizacao.Estado
	report.CodeState = int32(p.loja.Localizacao.CodigoEstado)
	report.City = p.loja.Localizacao.Cidade
	report.CodeCity = int32(p.loja.Localizacao.CodigoCidade)
	report.Neighborhood = p.loja.Localizacao.Bairro
	report.Partner = p.loja.Automacao
	report.Version = fmt.Sprintf("%g", p.loja.CodigoSistema)
}

// sendData envia o relatório em JSON para a API externa.
func (p *parseService) sendData(report models.Item) error {
	apiGateway := "https://api.painelmargem.com.br/gateway/report"
	jsonData, err := json.Marshal(report)
	if err != nil {
		return fmt.Errorf("erro ao converter para JSON: %w", err)
	}

	resp, err := http.Post(apiGateway, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao enviar POST para API: %w", err)
	}
	defer resp.Body.Close()
	p.logDetalhado("📡 [Parse] Resposta da API: %s", resp.Status)
	return nil
}

// logDetalhado imprime logs apenas se o ModoLog estiver definido como "detalhado".
func (p *parseService) logDetalhado(format string, args ...interface{}) {
	if p.modoLog == "detalhado" {
		log.Printf(format, args...)
	}
}
