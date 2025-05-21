package commands

import (
	"database/sql"
	"log"
	"time"

	"margem/robo/models/config"
	"margem/robo/repository"
)

// -----------------------------------------------------------------------------
// Configuração do horário extra (01:30 por padrão)
// -----------------------------------------------------------------------------
const extraHour = 1 // 01 h
const extraMin = 30 // 30 min

// Run mantém dois fluxos:
// 1) ticker normal para o dia corrente;
// 2) goroutine que dorme até o próximo HH:MM e executa o "extra".
func Run(
	reportRepos map[int]repository.ReportRepository,
	dbMap map[int]*sql.DB,
	interval int,
	cfg config.Config,
	modoLog string,
) {
	log.Printf("🚀 Run iniciado | ciclo normal: %d min | extra: %02d:%02d",
		interval, extraHour, extraMin)

	// Inicia o agendador extra em segundo‑plano.
	go scheduleExtra(reportRepos, dbMap, cfg, modoLog)

	// Loop principal (ciclo normal)
	ticker := time.NewTicker(time.Duration(interval) * time.Minute)
	defer ticker.Stop()

	// Contador para rotação do log (a cada 10 ciclos)
	cycleCount := 0

	for {
		executeRunCycle(reportRepos, dbMap, dateNow, cfg, modoLog)

		// Verifica e rotaciona o log a cada 10 ciclos
		cycleCount++
		if cycleCount >= 10 {
			checkAndRotateLog()
			cycleCount = 0
		}

		<-ticker.C
	}
}

// scheduleExtra aguarda até o próximo horário‑alvo, executa o extra
// e volta a aguardar para o dia seguinte.
func scheduleExtra(
	reportRepos map[int]repository.ReportRepository,
	dbMap map[int]*sql.DB,
	cfg config.Config,
	modoLog string,
) {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(),
			extraHour, extraMin, 0, 0, now.Location())

		if !next.After(now) { // se já passou hoje, agenda amanhã
			next = next.Add(24 * time.Hour)
		}

		sleep := time.Until(next)
		log.Printf("⏰ [Extra] Próxima atualização em %s (T‑%s)",
			next.Format("02/01 15:04"), sleep.Truncate(time.Second))

		time.Sleep(sleep)

		log.Println("⏰ [Extra] Iniciando sincronização dos 3 dias anteriores…")
		syncPreviousThreeDays(reportRepos, dbMap, cfg, modoLog, time.Now())

		// Verifica e rotaciona o log após o processamento extra
		checkAndRotateLog()
	}
}

// syncPreviousThreeDays dispara executeExtraCycle para D‑1, D‑2 e D‑3.
func syncPreviousThreeDays(
	reportRepos map[int]repository.ReportRepository,
	dbMap map[int]*sql.DB,
	cfg config.Config,
	modoLog string,
	now time.Time,
) {
	for i := 1; i <= 3; i++ {
		day := now.AddDate(0, 0, -i).Format("2006-01-02")
		ds := day + " 00:00:00"
		df := day + " 23:59:59"
		log.Printf("📆 [Extra] Sincronizando %s (dia ‑%d)", day, i)
		executeExtraCycle(reportRepos, dbMap, ds, df, cfg, modoLog)
	}
}

// checkAndRotateLog verifica e rotaciona o log se necessário
func checkAndRotateLog() {
	if err := reopenLogFile(); err != nil {
		log.Printf("⚠️ [Run] Erro ao rotacionar log: %v", err)
	}
}
