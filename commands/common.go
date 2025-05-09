package commands

import (
	"database/sql"
	"log"
	"margem/robo/models/config"
	"margem/robo/repository"
	"margem/robo/services" //  ←  novo import
	"time"
)

// -----------------------------------------------------------------------------
// Funções auxiliares e ciclos genéricos
// -----------------------------------------------------------------------------

// logSimpleSync emite um log compacto de sincronização bem‑sucedida.
func logSimpleSync(date string, duration time.Duration) {
	t, _ := time.Parse("2006-01-02", date[:10])
	log.Printf("✅ Sincronizado %s %s", t.Format("02/01/2006"), duration)
}

// executeRunCycle sincroniza o dia atual (ou ontem, dependendo de dateFunc).
func executeRunCycle(
	reportRepos map[int]repository.ReportRepository,
	dbMap map[int]*sql.DB,
	dateFunc func() (string, string),
	cfg config.Config,
	modoLog string,
) {
	log.Println("🔄 [Cycle] Iniciando ciclo de sincronização…")
	start := time.Now()

	ds, df := dateFunc()
	ok := 0
	for _, loja := range cfg.Lojas {
		repo, repoOK := reportRepos[loja.NumeroLoja]
		db, dbOK := dbMap[loja.NumeroLoja]
		if !repoOK || !dbOK || repo == nil || db == nil {
			log.Printf("⚠️ [Cycle] Loja %s sem repo/DB disponível.", loja.NomeFantasia)
			continue
		}

		service := services.NewParseService(loja, modoLog) //  ←  corrigido
		if err := service.Parse(ds, df); err != nil {
			log.Printf("❌ [Cycle] Erro na loja %s: %v", loja.NomeFantasia, err)
		} else {
			log.Printf("✅ [Cycle] Loja %s sincronizada (%s).", loja.NomeFantasia, ds[:10])
			ok++
		}
	}

	log.Printf("🔄 [Cycle] Concluído: %d/%d lojas em %s",
		ok, len(cfg.Lojas), time.Since(start))
}

// executeExtraCycle sincroniza um intervalo arbitrário (usado pelo “extra”).
func executeExtraCycle(
	reportRepos map[int]repository.ReportRepository,
	dbMap map[int]*sql.DB,
	ds, df string,
	cfg config.Config,
	modoLog string,
) {
	log.Printf("🔄 [Extra] Atualização extra %s → %s", ds, df)
	start := time.Now()
	ok := 0

	for _, loja := range cfg.Lojas {
		repo, repoOK := reportRepos[loja.NumeroLoja]
		db, dbOK := dbMap[loja.NumeroLoja]
		if !repoOK || !dbOK || repo == nil || db == nil {
			log.Printf("⚠️ [Extra] Loja %s sem repo/DB.", loja.NomeFantasia)
			continue
		}

		service := services.NewParseService(loja, modoLog) //  ←  corrigido
		if err := service.Parse(ds, df); err != nil {
			log.Printf("❌ [Extra] Erro na loja %s: %v", loja.NomeFantasia, err)
		} else {
			ok++
		}
	}
	log.Printf("🔄 [Extra] Finalizado: %d/%d lojas em %s",
		ok, len(cfg.Lojas), time.Since(start))
}

// dateNow devolve o range de hoje 00:00 → 23:59.
func dateNow() (string, string) {
	d := time.Now().Format("2006-01-02")
	return d + " 00:00:00", d + " 23:59:59"
}

// dateLastDay devolve o range de ontem 00:00 → 23:59.
func dateLastDay() (string, string) {
	d := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	return d + " 00:00:00", d + " 23:59:59"
}
