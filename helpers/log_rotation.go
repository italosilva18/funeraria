package helpers

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	// MaxLogSize é o tamanho máximo do arquivo de log em bytes (por exemplo, 10MB)
	MaxLogSize = 10 * 1024 * 1024
	// KeepPercentage é a porcentagem do log a ser mantida após a rotação
	KeepPercentage = 30
)

// RotateLogIfNeeded verifica se o arquivo de log excedeu o tamanho máximo
// e, se necessário, mantém apenas as entradas mais recentes (KeepPercentage)
func RotateLogIfNeeded(logFilePath string) error {
	fileInfo, err := os.Stat(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// O arquivo não existe, não precisa rotacionar
			return nil
		}
		return fmt.Errorf("erro ao verificar arquivo de log: %w", err)
	}

	// Verifica se o arquivo excedeu o tamanho máximo
	if fileInfo.Size() < MaxLogSize {
		return nil // Não precisa rotacionar
	}

	// Calcula a porcentagem de dados a manter
	bytesToKeep := int64(float64(fileInfo.Size()) * (KeepPercentage / 100.0))

	// Abre o arquivo de log
	file, err := os.Open(logFilePath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de log: %w", err)
	}

	// Move o cursor para a posição onde começam as entradas a serem mantidas
	_, err = file.Seek(fileInfo.Size()-bytesToKeep, io.SeekStart)
	if err != nil {
		file.Close()
		return fmt.Errorf("erro ao posicionar no arquivo de log: %w", err)
	}

	// Lê as entradas a serem mantidas
	buffer := make([]byte, bytesToKeep)
	bytesRead, err := file.Read(buffer)
	file.Close()
	if err != nil && err != io.EOF {
		return fmt.Errorf("erro ao ler arquivo de log: %w", err)
	}

	// Encontra o início da primeira linha completa
	// (para evitar começar com uma linha parcial)
	start := 0
	for i := 0; i < bytesRead; i++ {
		if buffer[i] == '\n' {
			start = i + 1
			break
		}
	}

	// Cria um arquivo temporário para as entradas a serem mantidas
	tempFile, err := os.CreateTemp("", "log-temp-*")
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo temporário: %w", err)
	}
	tempFilePath := tempFile.Name()

	// Escreve as entradas a serem mantidas no arquivo temporário
	_, err = tempFile.Write(buffer[start:bytesRead])
	tempFile.Close()
	if err != nil {
		os.Remove(tempFilePath)
		return fmt.Errorf("erro ao escrever no arquivo temporário: %w", err)
	}

	// Substitui o arquivo de log original pelo arquivo temporário
	err = os.Rename(tempFilePath, logFilePath)
	if err != nil {
		os.Remove(tempFilePath)
		return fmt.Errorf("erro ao substituir arquivo de log: %w", err)
	}

	log.Printf("✂️ Log rotacionado: mantidos %.2f%% dos logs mais recentes (%.2f MB)",
		KeepPercentage, float64(bytesRead-start)/(1024*1024))
	return nil
}

// CheckLogSize verifica o tamanho atual do arquivo de log
func CheckLogSize(logFilePath string) (int64, error) {
	fileInfo, err := os.Stat(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("erro ao verificar arquivo de log: %w", err)
	}
	return fileInfo.Size(), nil
}

// FormatFileSize converte bytes para uma representação mais legível
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
