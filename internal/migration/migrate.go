package migration

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"MysqlToMongo/internal/config"
	"MysqlToMongo/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

// setupLogging configura o log para arquivo e console
func setupLogging() (*os.File, error) {
	// Cria o diretório de logs se não existir
	if err := os.MkdirAll("tmp/logs", 0755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de logs: %v", err)
	}

	// Gera o nome do arquivo com timestamp
	timestamp := time.Now().Format("2006-01-02_15-04")
	logFile := filepath.Join("tmp", "logs", fmt.Sprintf("export_%s.log", timestamp))

	// Abre o arquivo de log
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo de log: %v", err)
	}

	// Configura o log para escrever tanto no arquivo quanto no console
	log.SetOutput(io.MultiWriter(os.Stdout, file))
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	return file, nil
}

// MigrateData executa o processo de migração dos dados
func MigrateData(config *config.Config, mysqlDB *sql.DB, mongoClient *mongo.Client) error {
	// Configura o logging
	logFile, err := setupLogging()
	if err != nil {
		return fmt.Errorf("erro ao configurar logging: %v", err)
	}
	defer logFile.Close()

	ctx := context.Background()
	collection := mongoClient.Database(config.MongoDB.Database).Collection(config.MongoDB.Collection)

	// Limpa a collection antes de começar
	log.Printf("Limpando collection '%s' existente...", config.MongoDB.Collection)
	if err := collection.Drop(ctx); err != nil {
		return fmt.Errorf("erro ao limpar collection: %v", err)
	}
	log.Printf("Collection '%s' limpa com sucesso!", config.MongoDB.Collection)
	log.Println("")

	// Obtém o total de registros (limitado a 5 milhões)
	var totalRecords int64
	err = mysqlDB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM (SELECT * FROM %s LIMIT 5000000) as t", config.MySQL.Table)).Scan(&totalRecords)
	if err != nil {
		return fmt.Errorf("erro ao contar registros: %v", err)
	}

	// Define o número de workers
	numWorkers := 5
	chunks := SplitWork(totalRecords, numWorkers)

	// Canais para controle
	errorChan := make(chan error, numWorkers)
	progressChan := make(chan int, numWorkers)
	var wg sync.WaitGroup

	// Inicia os workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		worker := &models.MigrationWorker{
			ID:           i + 1,
			StartID:      chunks[i].Start,
			EndID:        chunks[i].End,
			MySQLDB:      mysqlDB,
			MongoClient:  mongoClient,
			Config:       config,
			Wg:           &wg,
			ErrorChan:    errorChan,
			ProgressChan: progressChan,
		}

		go func(w *models.MigrationWorker) {
			defer w.Wg.Done()
			if err := w.ProcessBatch(ctx); err != nil {
				w.ErrorChan <- fmt.Errorf("erro no processador %d: %v", w.ID, err)
			}
		}(worker)
	}

	// Monitora o progresso
	go func() {
		totalProcessed := 0
		startTime := time.Now()
		reportThreshold := config.General.ReportThreshold
		isFirstReport := true

		for progress := range progressChan {
			totalProcessed += progress

			// Always show first report and when we reach the threshold
			if isFirstReport || totalProcessed >= reportThreshold {
				elapsed := time.Since(startTime)
				recordsPerSecond := float64(totalProcessed) / elapsed.Seconds()
				estimatedTotalTime := time.Duration(float64(totalRecords)/recordsPerSecond) * time.Second
				remainingTime := estimatedTotalTime - elapsed

				log.Printf("Progresso: %d/%d registros (%.2f%%) - Tempo decorrido: %v - Velocidade: %.2f registros/seg - Tempo restante estimado: %v",
					totalProcessed, totalRecords,
					float64(totalProcessed)/float64(totalRecords)*100,
					elapsed.Round(time.Second),
					recordsPerSecond,
					remainingTime.Round(time.Second))

				// Update the threshold for the next report
				reportThreshold = (totalProcessed/config.General.ReportThreshold + 1) * config.General.ReportThreshold
				isFirstReport = false
			}
		}

		// Show final progress after all processing is done
		elapsed := time.Since(startTime)
		recordsPerSecond := float64(totalProcessed) / elapsed.Seconds()
		log.Printf("Progresso: %d/%d registros (100.00%%) - Tempo total: %v - Velocidade média: %.2f registros/seg",
			totalProcessed, totalRecords,
			elapsed.Round(time.Second),
			recordsPerSecond)
	}()

	// Aguarda a conclusão
	wg.Wait()
	close(errorChan)
	close(progressChan)

	// Verifica erros
	for err := range errorChan {
		if err != nil {
			return err
		}
	}

	// Aguarda um momento para garantir que todas as operações foram concluídas
	time.Sleep(1 * time.Second)

	// Criar índices após a importação estar 100% completa
	log.Println("")
	log.Println("Criando índices...")
	if err := CreateIndexes(ctx, collection); err != nil {
		return fmt.Errorf("erro ao criar índices: %v", err)
	}

	return nil
}
