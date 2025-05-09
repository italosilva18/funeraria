package commands

import (
	"database/sql"
	"log"
	"margem/robo/models/config"
	"margem/robo/repository"
	"margem/robo/services"
	"time"
)

// TestRun executa o comando "simulation" para uma sincronização única do dia atual.
func TestRun(reportRepos map[int]repository.ReportRepository, dbMap map[int]*sql.DB, configData config.Config, modoLog string) {
	log.Println("🚀 [TestRun] Iniciando sincronização única (simulação) para o dia atual.")
	startTime := time.Now()
	ds, _ := dateNow() // ds: início, df: fim
	processedCount := 0

	for _, loja := range configData.Lojas {
		_, ok1 := reportRepos[loja.NumeroLoja]
		_, ok2 := dbMap[loja.NumeroLoja]
		if !ok1 || !ok2 {
			log.Printf("⚠️ [TestRun] Loja %s: Conexão ou repositório indisponível.", loja.NomeFantasia)
			continue
		}
		service := services.NewParseService(loja, modoLog)
		if err := service.Parse(ds, ds[:10]+" 23:59:59"); err != nil {
			log.Printf("❌ [TestRun] Erro na sincronização da loja %s: %v", loja.NomeFantasia, err)
		} else {
			log.Printf("✅ [TestRun] Loja '%s' (%d) -> data %s sincronizada OK.", loja.NomeFantasia, loja.NumeroLoja, ds[:10])
			processedCount++
		}
	}

	elapsed := time.Since(startTime)
	log.Printf("⏱ [TestRun] Finalizado. Lojas sincronizadas: %d/%d. Tempo total: %s", processedCount, len(configData.Lojas), elapsed)
}
