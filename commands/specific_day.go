package commands

import (
	"database/sql"
	"log"
	"margem/robo/models/config"
	"margem/robo/repository"
	"margem/robo/services"
	"time"
)

// SpecificDay executa uma sincronização única para uma data específica.
func SpecificDay(reportRepos map[int]repository.ReportRepository, dbMap map[int]*sql.DB, dateStr string, configData config.Config, modoLog string) {
	log.Printf("📆 [SpecificDay] Iniciando sincronização para a data: %s", dateStr)
	ds := dateStr + " 00:00:00"
	df := dateStr + " 23:59:59"
	startTime := time.Now()
	processedCount := 0

	for _, loja := range configData.Lojas {
		_, ok1 := reportRepos[loja.NumeroLoja]
		_, ok2 := dbMap[loja.NumeroLoja]
		if !ok1 || !ok2 {
			log.Printf("⚠️ [SpecificDay] Loja %s: Conexão ou repositório indisponível.", loja.NomeFantasia)
			continue
		}
		service := services.NewParseService(loja, modoLog)
		if err := service.Parse(ds, df); err != nil {
			log.Printf("❌ [SpecificDay] Erro na sincronização da loja %s: %v", loja.NomeFantasia, err)
		} else {
			log.Printf("✅ [SpecificDay] Loja '%s' sincronizada com sucesso para %s.", loja.NomeFantasia, dateStr)
			processedCount++
		}
	}
	elapsed := time.Since(startTime)
	log.Printf("⏱ [SpecificDay] Finalizado. Lojas sincronizadas: %d/%d. Tempo total: %s", processedCount, len(configData.Lojas), elapsed)
}
