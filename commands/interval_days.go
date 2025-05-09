package commands

import (
	"database/sql"
	"log"
	"margem/robo/models/config"
	"margem/robo/repository"
	"margem/robo/services"
	"time"
)

// IntervalDays executa o comando "interval-days" para um intervalo de datas
func IntervalDays(reportRepos map[int]repository.ReportRepository, dbMap map[int]*sql.DB, startDate, endDate string, configData config.Config, modoLog string) {
	log.Printf("-> IntervalDays: %s a %s", startDate, endDate)
	layout := "2006-01-02"

	sTime, err := time.Parse(layout, startDate)
	if err != nil {
		log.Printf("   Erro parse startDate: %v", err)
		return
	}
	eTime, err := time.Parse(layout, endDate)
	if err != nil {
		log.Printf("   Erro parse endDate: %v", err)
		return
	}
	if eTime.Before(sTime) {
		log.Printf("   Data final %s < data inicial %s", endDate, startDate)
		return
	}

	for d := sTime; !d.After(eTime); d = d.AddDate(0, 0, 1) {
		ds := d.Format(layout) + " 00:00:00"
		df := d.Format(layout) + " 23:59:59"
		startDay := time.Now()

		processedCount := 0
		for _, loja := range configData.Lojas {
			_, ok1 := reportRepos[loja.NumeroLoja]
			_, ok2 := dbMap[loja.NumeroLoja]
			if !ok1 || !ok2 {
				continue
			}
			service := services.NewParseService(loja, modoLog)
			if err := service.Parse(ds, df); err != nil {
				continue
			} else {
				processedCount++
			}
		}
		if processedCount == len(configData.Lojas) {
			elapsed := time.Since(startDay)
			logSimpleSync(ds, elapsed)
		}
	}
}
