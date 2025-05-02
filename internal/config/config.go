package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config representa a configuração geral da aplicação
type Config struct {
	MySQL   MySQLConfig    `json:"mysql"`
	MongoDB MongoDBConfig  `json:"mongodb"`
	General GeneralConfig  `json:"general"`
	Mapping *MappingConfig `json:"-"` // Não será carregado do config.json
}

// MySQLConfig representa a configuração de conexão com o MySQL
type MySQLConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	Table    string `json:"table"`
}

// MongoDBConfig representa a configuração de conexão com o MongoDB
type MongoDBConfig struct {
	URI        string `json:"uri"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
}

// GeneralConfig representa configurações gerais da aplicação
type GeneralConfig struct {
	BatchSize       int `json:"batch_size"`
	NumWorkers      int `json:"num_workers"`
	ReportThreshold int `json:"report_threshold"`
}

// MappingConfig representa o mapeamento das colunas
type MappingConfig struct {
	Pessoas struct {
		CPF             int `json:"cpf"`
		Nome            int `json:"nome"`
		Nasc            int `json:"nasc"`
		Renda           int `json:"renda"`
		AffinityScore   int `json:"affinity_score"`
		AffinityPercent int `json:"affinity_percent"`
		Sexo            int `json:"sexo"`
		CBO             int `json:"cbo"`
		Mae             int `json:"mae"`
		Nota            int `json:"nota"`
		Banco           int `json:"banco"`
		CPFConjuge      int `json:"cpf_conjuge"`
		ServPublico     int `json:"serv_publico"`
		DataObito       int `json:"data_obito"`
		Cidade          int `json:"cidade"`
		Endereco        int `json:"endereco"`
		Bairro          int `json:"bairro"`
		CEP             int `json:"cep"`
		UF              int `json:"uf"`
		DataAtualizacao int `json:"data_atualizacao"`
		Contatos        struct {
			Telefones []int `json:"telefones"`
			Emails    []int `json:"emails"`
		} `json:"contatos"`
	} `json:"pessoas"`
}

// LoadConfig carrega a configuração do arquivo config.json
func LoadConfig() (*Config, error) {
	// Carrega config.json
	configPath := filepath.Join("config", "config.json")
	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	// Carrega mapping.json
	mappingPath := filepath.Join("config", "mapping.json")
	mappingFile, err := os.ReadFile(mappingPath)
	if err != nil {
		return nil, err
	}

	var mapping MappingConfig
	if err := json.Unmarshal(mappingFile, &mapping); err != nil {
		return nil, err
	}

	config.Mapping = &mapping
	return &config, nil
}
