package migration

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Função para criar índices
func CreateIndexes(ctx context.Context, collection *mongo.Collection) error {
	// Configuração da collation case-sensitive apenas para o nome
	nomeCollation := options.Collation{
		Locale:   "pt",
		Strength: 3,
	}

	// Índice único para CPF
	cpfIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "cpf", Value: 1}},
		Options: options.Index().
			SetUnique(true),
	}

	// Índice para nome (case-sensitive)
	nomeIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "nome", Value: 1}},
		Options: options.Index().
			SetCollation(&nomeCollation),
	}

	// Índice para emails (dentro do array de contatos)
	emailIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "contatos.emails", Value: 1}},
	}

	// Índice para telefones (dentro do array de contatos)
	telefoneIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "contatos.telefones", Value: 1}},
	}

	// Criar todos os índices
	indexes := []mongo.IndexModel{cpfIndex, nomeIndex, emailIndex, telefoneIndex}
	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("erro ao criar índices: %v", err)
	}

	log.Println("Índices criados com sucesso!")
	return nil
}
