package helpers

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

// DateTo1AAMMDD converte uma data no formato "2006-01-02 15:04:05"
// para uma string no formato "1AAMMDD".
// Exemplo: "2025-04-06 00:00:00" => "1250406"
func DateTo1AAMMDD(dateStr string) (string, error) {
	date, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return "", err
	}

	// Extrai ano, mês e dia
	year, month, day := date.Date()

	// Converte ano para formato "AA" (pegando os 2 últimos dígitos)
	yearStr := strconv.Itoa(year)[2:] // ex: 2025 -> "25"

	// Constrói "1AAMMDD"
	resultStr := "1" + yearStr + fmt.Sprintf("%02d%02d", month, day)
	return resultStr, nil
}

// GetLastNDays retorna uma lista de strings com as datas dos
// últimos `n` dias (no formato YYYY-MM-DD), usando o timezone `loc`.
//
// Exemplo: Se hoje é 07/04/2025, GetLastNDays(3, loc)
// retorna: ["2025-04-04", "2025-04-05", "2025-04-06"].
func GetLastNDays(n int, loc *time.Location) []string {
	dates := make([]string, 0, n)
	now := time.Now().In(loc)

	for i := n; i >= 1; i-- {
		day := now.AddDate(0, 0, -i)
		dates = append(dates, day.Format("2006-01-02"))
	}
	return dates
}

// GetTimezone retorna a localização (timezone) usada no sistema.
// Por padrão, tenta "America/Sao_Paulo". Caso não encontre,
// cai para time.UTC e imprime um aviso no log.
func GetTimezone() *time.Location {
	tz := "America/Sao_Paulo"
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Printf("Aviso: Fuso horário '%s' não encontrado, usando UTC. Erro: %v", tz, err)
		return time.UTC
	}
	return loc
}
