package migration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"MysqlToMongo/internal/config"
	"MysqlToMongo/internal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

// MigrateData executa o processo de migração dos dados
func MigrateData(config *config.Config, mysqlDB *sql.DB, mongoClient *mongo.Client) error {
	ctx := context.Background()
	collection := mongoClient.Database(config.MongoDB.Database).Collection(config.MongoDB.Collection)

	// Limpa a collection antes de começar
	log.Println("Limpando collection existente...")
	if err := collection.Drop(ctx); err != nil {
		return fmt.Errorf("erro ao limpar collection: %v", err)
	}
	log.Println("Collection limpa com sucesso!")
	log.Println("")

	// Obtém o total de registros
	var totalRecords int64
	err := mysqlDB.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", config.MySQL.Table)).Scan(&totalRecords)
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
		reportThreshold := 1000000 // Report every 1 million records

		for progress := range progressChan {
			totalProcessed += progress

			// Only report when we've processed another million records
			if totalProcessed >= reportThreshold {
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

				// Update the threshold for the next million
				reportThreshold = (totalProcessed/1000000 + 1) * 1000000
			}
		}
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
