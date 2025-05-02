package models

import (
	"context"
	"fmt"
	"runtime"

	"MysqlToMongo/internal/converter"
)

// Função para calcular o limite de memória
func calculateMemoryLimit() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	totalMemory := m.Sys
	memoryLimit := int64(float64(totalMemory) * 0.7)
	return memoryLimit
}

// ProcessBatch processa um lote de registros
func (w *MigrationWorker) ProcessBatch(ctx context.Context) error {
	collection := w.MongoClient.Database(w.Config.MongoDB.Database).Collection(w.Config.MongoDB.Collection)

	// Query para obter apenas os registros do worker usando LIMIT e OFFSET
	query := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ?", w.Config.MySQL.Table)
	rows, err := w.MySQLDB.Query(query, w.EndID-w.StartID+1, w.StartID-1)
	if err != nil {
		return fmt.Errorf("erro na consulta MySQL: %v", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("erro ao obter colunas: %v", err)
	}

	// Prepare slice for values
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// Calcula o tamanho do lote baseado na memória disponível
	memoryLimit := calculateMemoryLimit()
	estimatedDocSize := int64(1024)                        // Estimativa de 1KB por documento
	batchSize := int(memoryLimit / (estimatedDocSize * 5)) // Divide por 5 workers
	if batchSize > w.Config.General.BatchSize {
		batchSize = w.Config.General.BatchSize
	}

	batch := make([]interface{}, 0, batchSize)
	count := 0

	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			return fmt.Errorf("erro ao escanear linha: %v", err)
		}

		// Create document for MongoDB
		doc := OrderedDocument{}

		// Map pessoa fields
		p := w.Config.Mapping.Pessoas
		doc.CPF = converter.ConvertBinaryToString(values[p.CPF-1])
		doc.Nome = converter.ConvertBinaryToString(values[p.Nome-1])
		doc.Nasc = converter.ConvertToDatePtr(values[p.Nasc-1])
		doc.Renda = converter.ConvertToDecimal(values[p.Renda-1])
		doc.AffinityScore = converter.ConvertToDecimal(values[p.AffinityScore-1])
		doc.AffinityPercent = converter.ConvertToDecimal(values[p.AffinityPercent-1])
		doc.Sexo = converter.ConvertBinaryToString(values[p.Sexo-1])
		doc.CBO = converter.ConvertBinaryToString(values[p.CBO-1])
		doc.Mae = converter.ConvertBinaryToString(values[p.Mae-1])
		doc.Nota = converter.ConvertBinaryToString(values[p.Nota-1])
		doc.Banco = converter.ConvertBinaryToString(values[p.Banco-1])
		doc.CPFConjuge = converter.ConvertOptionalField(values[p.CPFConjuge-1])
		doc.ServPublico = converter.ConvertOptionalField(values[p.ServPublico-1])
		doc.DataObito = converter.ConvertToDatePtr(values[p.DataObito-1])
		doc.Cidade = converter.ConvertBinaryToString(values[p.Cidade-1])
		doc.Endereco = converter.ConvertBinaryToString(values[p.Endereco-1])
		doc.Bairro = converter.ConvertOptionalField(values[p.Bairro-1])
		doc.CEP = converter.ConvertBinaryToString(values[p.CEP-1])
		doc.UF = converter.ConvertBinaryToString(values[p.UF-1])
		doc.DataAtualizacao = converter.ConvertToTimePtr(values[p.DataAtualizacao-1])

		// Map telefones
		telefones := make([]interface{}, 0)
		for _, pos := range p.Contatos.Telefones {
			if values[pos-1] != nil {
				telefone := converter.ConvertBinaryToString(values[pos-1])
				if str, ok := telefone.(string); ok && str != "" {
					telefones = append(telefones, str)
				}
			}
		}
		doc.Contatos.Telefones = telefones

		// Map emails
		emails := make([]interface{}, 0)
		for _, pos := range p.Contatos.Emails {
			if values[pos-1] != nil {
				email := converter.ConvertBinaryToString(values[pos-1])
				if str, ok := email.(string); ok && str != "" {
					emails = append(emails, str)
				}
			}
		}
		doc.Contatos.Emails = emails

		batch = append(batch, doc)
		count++

		// Insert batch when it reaches the batch size
		if len(batch) >= batchSize {
			if _, err := collection.InsertMany(ctx, batch); err != nil {
				return fmt.Errorf("erro ao inserir lote: %v", err)
			}
			w.ProgressChan <- len(batch)
			batch = batch[:0]
		}
	}

	// Insert remaining documents
	if len(batch) > 0 {
		if _, err := collection.InsertMany(ctx, batch); err != nil {
			return fmt.Errorf("erro ao inserir lote final: %v", err)
		}
		w.ProgressChan <- len(batch)
	}

	return nil
}
