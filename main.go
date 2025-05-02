package main

import (
	"context"
	"log"
	"time"

	"MysqlToMongo/internal/config"
	"MysqlToMongo/internal/database"
	"MysqlToMongo/internal/migration"
)

func main() {
	// Inicia o timer
	startTime := time.Now()

	// Carrega configuração
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	// Conecta ao MySQL
	mysqlDB, err := database.ConnectMySQL(config)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MySQL: %v", err)
	}
	defer mysqlDB.Close()

	// Conecta ao MongoDB
	mongoClient, err := database.ConnectMongoDB(config)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Inicia migração
	log.Println("Iniciando migração...")
	if err := migration.MigrateData(config, mysqlDB, mongoClient); err != nil {
		log.Fatalf("Erro durante a migração: %v", err)
	}

	// Calcula e mostra o tempo total
	duration := time.Since(startTime)
	log.Printf("Migração concluída com sucesso em %v!", duration)
}
