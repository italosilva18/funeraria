package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	// Embutir dados de fuso horário no binário
	_ "time/tzdata"

	// Pacotes do projeto
	"margem/robo/commands"
	"margem/robo/connections"
	"margem/robo/helpers"
	"margem/robo/models/config"
	"margem/robo/querys"
	"margem/robo/querys/concinco"
	"margem/robo/querys/derevo"
	"margem/robo/querys/gestor"
	"margem/robo/querys/rms"
	"margem/robo/repository"

	// Drivers de banco de dados
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/nakagami/firebirdsql"
	_ "github.com/sijms/go-ora/v2"

	"github.com/urfave/cli/v2"
)

var (
	// ConfigData armazena as configurações carregadas do arquivo JSON
	ConfigData config.Config
	// ModoLog define o nível do log ("detalhado" ou "resumido")
	ModoLog string
)

// containsHelp verifica se algum dos argumentos é "--help" ou "-h"
func containsHelp(args []string) bool {
	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

// printUsageGuide exibe um guia rápido de uso, apenas quando solicitado.
func printUsageGuide() {
	guide := `
CLI-Robo - Guia de Uso:

Global Options:
   --log <valor>        Define o modo de log: "detalhado" para logs completos ou "resumido" (padrão)
   --help, -h           Exibe esta mensagem de ajuda
   --version, -v        Exibe a versão da aplicação

Comandos:
   run                  Executa o robô em loop para sincronizar os dados do dia atual.
                        Exemplo: robo.exe run --interval 30

   last-day             Executa o robô em loop para sincronizar os dados do dia anterior.
                        Exemplo: robo.exe last-day --interval 30

   test                 Testa a conexão com os bancos de dados.
                        Exemplo: robo.exe test

   simulation           Executa uma sincronização única para o dia atual (simulação).
                        Exemplo: robo.exe simulation

   specific-day         Sincroniza uma data específica (YYYY-MM-DD).
                        Exemplo: robo.exe specific-day --date 2025-04-06

   interval-days        Sincroniza um intervalo de datas.
                        Exemplo: robo.exe interval-days --dateStart 2025-04-01 --dateEnd 2025-04-05
`
	fmt.Println(guide)
}

func main() {
	// Configura o formato dos logs
	log.SetFlags(log.Ldate | log.Ltime)

	// Configura o arquivo de log
	fmt.Println("⚙️ Configurando arquivo de Log em ./config/log.txt")
	logFilePath := "./config/log.txt"

	// Garante a existência do arquivo e as permissões corretas
	if err := helpers.EnsureLogFile(logFilePath); err != nil {
		log.Fatalf("Erro ao garantir arquivo de log: %v", err)
	}

	// Verifica o tamanho atual do log antes da rotação
	currentSize, err := helpers.CheckLogSize(logFilePath)
	if err == nil && currentSize > 0 {
		fmt.Printf("📊 Tamanho atual do log: %s\n", helpers.FormatFileSize(currentSize))
	}

	// Rotaciona o log se necessário antes de abri-lo para escrita
	if err := helpers.RotateLogIfNeeded(logFilePath); err != nil {
		fmt.Printf("⚠️ Aviso: Erro ao rotacionar log: %v\n", err)
	}

	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	logFile, err := os.OpenFile(logFilePath, flags, 0666)
	if err != nil {
		log.Fatalf("Erro fatal ao abrir arquivo de log: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	commands.SetLogFile(logFile, logFilePath)

	// Exibe o guia de uso somente se solicitado
	if len(os.Args) < 2 || containsHelp(os.Args[1:]) {
		printUsageGuide()
	}

	log.Println("=== CLI-Robo iniciado ===")
	log.Println("Versão: 2.2.0")

	cliApp := &cli.App{
		Name:     "CLI-Robo",
		Usage:    "Robô de integração e sincronização de dados de lojas.",
		Version:  "2.2.0",
		Compiled: time.Now(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log",
				Usage:       "Define o modo de log (detalhado ou resumido)",
				Value:       "resumido",
				Destination: &ModoLog,
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Executa o robô em loop para sincronizar os dados do dia atual.",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "interval",
						Usage:   "Intervalo de execução em minutos",
						Value:   30,
						Aliases: []string{"i"},
					},
				},
				Action: func(ctx *cli.Context) error {
					log.Println("==> [run] Iniciado: Sincronização contínua para o dia atual.")
					reportRepos, dbMap, err := Teste()
					if err != nil {
						log.Printf("❌ [run] Erro crítico na função Teste: %v", err)
						return err
					}
					defer closeConnections(dbMap)
					interval := ctx.Int("interval")
					log.Printf("🕒 [run] Intervalo configurado: %d minutos.", interval)
					commands.Run(reportRepos, dbMap, interval, ConfigData, ModoLog)
					log.Println("✅ [run] Finalizado.")
					return nil
				},
			},
			{
				Name:  "last-day",
				Usage: "Executa o robô em loop para sincronizar os dados do dia anterior.",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "interval",
						Usage:   "Intervalo de execução em minutos",
						Value:   30,
						Aliases: []string{"i"},
					},
				},
				Action: func(ctx *cli.Context) error {
					log.Println("==> [last-day] Iniciado: Sincronização contínua para o dia anterior.")
					reportRepos, dbMap, err := Teste()
					if err != nil {
						log.Printf("❌ [last-day] Erro crítico na função Teste: %v", err)
						return err
					}
					defer closeConnections(dbMap)
					interval := ctx.Int("interval")
					log.Printf("🕒 [last-day] Intervalo configurado: %d minutos.", interval)
					commands.LastDay(reportRepos, dbMap, interval, ConfigData, ModoLog)
					log.Println("✅ [last-day] Finalizado.")
					return nil
				},
			},
			{
				Name:  "test",
				Usage: "Testa a conexão com os bancos de dados.",
				Action: func(ctx *cli.Context) error {
					log.Println("🔍 [test] Iniciado: Testando conexões com os bancos.")
					_, dbMap, err := Teste()
					if err != nil {
						log.Printf("❌ [test] Erro durante o teste de conexão: %v", err)
						return err
					}
					defer closeConnections(dbMap)
					log.Println("✅ [test] Conexões testadas com sucesso.")
					return nil
				},
			},
			{
				Name:  "simulation",
				Usage: "Executa uma sincronização única para o dia atual (simulação).",
				Action: func(ctx *cli.Context) error {
					log.Println("🚀 [simulation] Iniciado: Sincronização única para o dia atual.")
					reportRepos, dbMap, err := Teste()
					if err != nil {
						log.Printf("❌ [simulation] Erro crítico na função Teste: %v", err)
						return err
					}
					defer closeConnections(dbMap)
					commands.TestRun(reportRepos, dbMap, ConfigData, ModoLog)
					log.Println("✅ [simulation] Finalizado.")
					return nil
				},
			},
			{
				Name:  "specific-day",
				Usage: "Executa uma sincronização única para uma data específica (YYYY-MM-DD).",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "date",
						Usage:    "Data específica no formato YYYY-MM-DD (ex: 2025-04-06)",
						Required: true,
						Aliases:  []string{"d"},
					},
				},
				Action: func(ctx *cli.Context) error {
					dateStr := ctx.String("date")
					log.Printf("📆 [specific-day] Iniciado para a data: %s", dateStr)
					_, err := time.Parse("2006-01-02", dateStr)
					if err != nil {
						log.Printf("❌ [specific-day] Data inválida. Formato correto: YYYY-MM-DD. Erro: %v", err)
						return fmt.Errorf("formato de data inválido: %s", dateStr)
					}
					reportRepos, dbMap, err := Teste()
					if err != nil {
						log.Printf("❌ [specific-day] Erro crítico na função Teste: %v", err)
						return err
					}
					defer closeConnections(dbMap)
					commands.SpecificDay(reportRepos, dbMap, dateStr, ConfigData, ModoLog)
					log.Println("✅ [specific-day] Finalizado.")
					return nil
				},
			},
			{
				Name:  "interval-days",
				Usage: "Executa sincronização para um intervalo de datas [start, end].",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "dateStart",
						Usage:    "Data de início no formato YYYY-MM-DD",
						Required: true,
						Aliases:  []string{"s"},
					},
					&cli.StringFlag{
						Name:     "dateEnd",
						Usage:    "Data de fim no formato YYYY-MM-DD",
						Required: true,
						Aliases:  []string{"e"},
					},
				},
				Action: func(ctx *cli.Context) error {
					dateStartStr := ctx.String("dateStart")
					dateEndStr := ctx.String("dateEnd")
					log.Printf("📆 [interval-days] Iniciado para o intervalo: %s a %s", dateStartStr, dateEndStr)
					_, err1 := time.Parse("2006-01-02", dateStartStr)
					_, err2 := time.Parse("2006-01-02", dateEndStr)
					if err1 != nil || err2 != nil {
						log.Printf("❌ [interval-days] Datas inválidas. Início: %v, Fim: %v", err1, err2)
						return fmt.Errorf("datas inválidas: início '%s', fim '%s'", dateStartStr, dateEndStr)
					}
					reportRepos, dbMap, err := Teste()
					if err != nil {
						log.Printf("❌ [interval-days] Erro crítico na função Teste: %v", err)
						return err
					}
					defer closeConnections(dbMap)
					commands.IntervalDays(reportRepos, dbMap, dateStartStr, dateEndStr, ConfigData, ModoLog)
					log.Println("✅ [interval-days] Finalizado.")
					return nil
				},
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Printf("❌ Erro final ao executar CLI: %v", err)
		log.Println("=== Aplicação finalizada com erro ===")
		os.Exit(1)
	}
	log.Println("🏁 Aplicação finalizada normalmente.")
}

// Teste lê a configuração, estabelece conexões e cria repositórios.
func Teste() (map[int]repository.ReportRepository, map[int]*sql.DB, error) {
	fmt.Println("🛠️ Lendo configurações do arquivo: config/configuracao.json")
	log.Println("-> Teste: Carregando configuração...")

	cfg, err := loadConfig("./config/configuracao.json")
	if err != nil {
		log.Printf("❌ Erro ao carregar config: %v", err)
		return nil, nil, err
	}
	ConfigData = *cfg
	log.Printf("✅ Configuração carregada com sucesso. %d loja(s) encontrada(s).", len(ConfigData.Lojas))

	reportRepos := make(map[int]repository.ReportRepository)
	dbMap := make(map[int]*sql.DB)
	startTime := time.Now()
	log.Println("-> Teste: Iniciando teste de conexão e criação de repositórios para as lojas...")

	for _, loja := range ConfigData.Lojas {
		log.Printf("   ---\n   Loja: %s (Número: %d, Automação: %s, Sistema: %.1f)",
			loja.NomeFantasia, loja.NumeroLoja, loja.Automacao, loja.CodigoSistema)

		dbHandler, err := connections.GetConnection(loja.Automacao, loja)
		if err != nil {
			log.Printf("      ❌ Erro ao obter handler para a loja %s: %v", loja.NomeFantasia, err)
			continue
		}

		db, err := dbHandler.DB()
		if err != nil {
			log.Printf("      ❌ Erro ao abrir DB para a loja %s: %v", loja.NomeFantasia, err)
			continue
		}

		if err = db.Ping(); err != nil {
			log.Printf("      ❌ Ping falhou para a loja %s: %v", loja.NomeFantasia, err)
			db.Close()
			continue
		}
		log.Printf("💾 Conexão com banco OK para a loja %s.", loja.NomeFantasia)

		qh, err := getQueryHandlerBySystem(loja)
		if err != nil {
			log.Printf("      ❌ Erro obtendo QueryHandler para a loja %s: %v", loja.NomeFantasia, err)
			db.Close()
			continue
		}
		log.Printf("⚙️ Automação '%s' (Sistema %.1f) carregada com sucesso.", loja.Automacao, loja.CodigoSistema)

		var repo repository.ReportRepository
		switch strings.ToUpper(strings.TrimSpace(loja.Automacao)) {
		case "RMS":
			repo = repository.NewReportRepository(db, qh)
		case "GESTOR":
			repo = repository.NewReportGestorRepository(db, qh)
		case "DEREVO":
			repo = repository.NewReportDerevoRepository(db, qh)
		case "ARIUS":
			repo = repository.NewReportRepository(db, qh)
		case "ARIUS-ERP":
			repo = repository.NewReportRepository(db, qh)
		case "CONCINCO":
			repo = repository.NewReportConcnicoRepository(db, qh)
		default:
			log.Printf("      ❌ Automação '%s' não suportada para a loja %s.", loja.Automacao, loja.NomeFantasia)
			db.Close()
			continue
		}
		log.Printf("🏚 Repositório criado para a loja %s.", loja.NomeFantasia)
		reportRepos[loja.NumeroLoja] = repo
		dbMap[loja.NumeroLoja] = db
	}

	elapsed := time.Since(startTime)
	log.Printf("⏱ Teste finalizado. Lojas inicializadas: %d em %s.", len(reportRepos), elapsed)
	if len(reportRepos) == 0 && len(ConfigData.Lojas) > 0 {
		log.Println("⚠️ Aviso: Nenhuma loja pôde ser inicializada com sucesso.")
	}

	return reportRepos, dbMap, nil
}

func getQueryHandlerBySystem(loja config.Loja) (querys.QueryHandler, error) {
	automacao := strings.ToUpper(strings.TrimSpace(loja.Automacao))
	log.Printf("🔧 Automação normalizada: '%s', Código do Sistema: %.1f", automacao, loja.CodigoSistema)
	switch automacao {
	case "RMS":
		versionStr := fmt.Sprintf("%.1f", loja.CodigoSistema)
		switch versionStr {
		case "51.1":
			return &rms.Rms511{}, nil
		case "51.2":
			return &rms.Rms512{}, nil
		}
		return nil, fmt.Errorf("❌ Versão RMS '%g' não implementada", loja.CodigoSistema)
	case "GESTOR":
		versionStr := fmt.Sprintf("%.1f", loja.CodigoSistema)
		switch versionStr {
		case "48.1":
			return &gestor.Gestor481{}, nil
		case "48.2":
			return &gestor.Gestor482{}, nil
		}
		return nil, fmt.Errorf("❌ Versão GESTOR '%g' não implementada", loja.CodigoSistema)
	case "DEREVO":
		versionStr := fmt.Sprintf("%.1f", loja.CodigoSistema)
		switch versionStr {
		case "52.1":
			return &derevo.Derevo521{}, nil
		}
		return nil, fmt.Errorf("❌ Versão DEREVO '%g' não implementada", loja.CodigoSistema)
	case "CONCINCO":
		versionStr := fmt.Sprintf("%.1f", loja.CodigoSistema)
		switch versionStr {
		case "28.1":
			return &concinco.Concinco281{}, nil
		}
		return nil, fmt.Errorf("❌ Versão CONCINCO '%g' não implementada", loja.CodigoSistema)
	default:
		return nil, fmt.Errorf("❌ Automação '%s' não reconhecida", automacao)
	}
}

func loadConfig(filePath string) (*config.Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler config '%s': %w", filePath, err)
	}
	var cfg config.Config
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("erro ao interpretar JSON '%s': %w", filePath, err)
	}
	return &cfg, nil
}

func closeConnections(dbMap map[int]*sql.DB) {
	log.Println("🔌 Fechando conexões...")
	closedCount := 0
	for lojaNum, db := range dbMap {
		if db != nil {
			log.Printf("   Fechando conexão da loja %d...", lojaNum)
			err := db.Close()
			if err != nil {
				log.Printf("   ⚠️ Erro ao fechar conexão da loja %d: %v", lojaNum, err)
			} else {
				log.Printf("   ✅ Conexão da loja %d fechada.", lojaNum)
				closedCount++
			}
		}
	}
	log.Printf("🔒 Conexões fechadas: %d.", closedCount)
}
