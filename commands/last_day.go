package commands

import (
	"database/sql"
	"log"
	"margem/robo/models/config"
	"margem/robo/repository"
	"time"
)

// LastDay executa o comando "last-day" em loop para sincronizar os dados do dia anterior.
func LastDay(reportRepos map[int]repository.ReportRepository, dbMap map[int]*sql.DB, interval int, configData config.Config, modoLog string) {
	log.Printf("📆 [LastDay] Iniciando sincronização contínua para o dia anterior (intervalo: %d min).", interval)
	ticker := time.NewTicker(time.Duration(interval) * time.Minute)
	defer ticker.Stop()

	executeRunCycle(reportRepos, dbMap, dateLastDay, configData, modoLog)
	for range ticker.C {
		executeRunCycle(reportRepos, dbMap, dateLastDay, configData, modoLog)
	}
}
