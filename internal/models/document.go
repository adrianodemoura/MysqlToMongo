package models

import (
	"sync"

	"MysqlToMongo/internal/config"
	"database/sql"

	"go.mongodb.org/mongo-driver/mongo"
)

// OrderedDocument representa a estrutura ordenada do documento no MongoDB
type OrderedDocument struct {
	CPF             any `bson:"cpf"`
	Nome            any `bson:"nome"`
	Nasc            any `bson:"nasc"`
	Renda           any `bson:"renda"`
	AffinityScore   any `bson:"affinity_score"`
	AffinityPercent any `bson:"affinity_percent"`
	Sexo            any `bson:"sexo"`
	CBO             any `bson:"cbo"`
	Mae             any `bson:"mae"`
	Nota            any `bson:"nota"`
	Banco           any `bson:"banco"`
	CPFConjuge      any `bson:"cpf_conjuge"`
	ServPublico     any `bson:"serv_publico"`
	DataObito       any `bson:"data_obito"`
	Cidade          any `bson:"cidade"`
	Endereco        any `bson:"endereco"`
	Bairro          any `bson:"bairro"`
	CEP             any `bson:"cep"`
	UF              any `bson:"uf"`
	DataAtualizacao any `bson:"data_atualizacao"`
	Contatos        struct {
		Telefones []any `bson:"telefones"`
		Emails    []any `bson:"emails"`
	} `bson:"contatos"`
}

// MigrationWorker representa um worker para processamento paralelo
type MigrationWorker struct {
	ID           int
	StartID      int64
	EndID        int64
	MySQLDB      *sql.DB
	MongoClient  *mongo.Client
	Config       *config.Config
	Wg           *sync.WaitGroup
	ErrorChan    chan error
	ProgressChan chan int
}
