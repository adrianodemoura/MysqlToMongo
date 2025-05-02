# MySQL to MongoDB Migration Tool

Este é um script em Go para migrar dados de uma tabela MySQL para uma coleção MongoDB, com suporte a processamento paralelo e conversão automática de tipos de dados.

## Estrutura do Projeto

```
MysqlToMongo/
├── config/
│   ├── config.json     # Configurações de conexão e parâmetros gerais
│   └── mapping.json    # Mapeamento estático das colunas
├── internal/
│   ├── config/         # Gerenciamento de configurações
│   ├── converter/      # Funções de conversão de tipos
│   ├── database/       # Conexões com bancos de dados
│   ├── migration/      # Lógica de migração
│   └── models/         # Estruturas de dados
├── scripts/            # Scripts utilitários
│   ├── buscaTelefone.sh    # Script para busca de telefones
│   ├── corrigirCPFs.sh     # Script para correção de CPFs
│   └── outros scripts...
├── tmp/
│   └── logs/          # Diretório para arquivos de log
├── main.go            # Ponto de entrada da aplicação
└── README.md          # Este arquivo
```

## Configuração

### config.json
```json
{
    "mysql": {
        "host": "localhost",
        "port": 3306,
        "user": "",
        "password": "",
        "database": "",
        "table": ""
    },
    "mongodb": {
        "uri": "mongodb://userMongo:senhaUserMongo@localhost:27017/database",
        "database": "",
        "collection": ""
    },
    "general": {
        "batch_size": 1000,
        "num_workers": 5
    }
}
```

### mapping.json
```json
{
    "pessoas": {
        "cpf": 1,
        "nome": 2,
        "nasc": 3,
        "renda": 4,
        "affinity_score": 5,
        "affinity_percent": 6,
        "sexo": 7,
        "cbo": 8,
        "mae": 9,
        "nota": 10,
        "banco": 11,
        "cpf_conjuge": 12,
        "serv_publico": 13,
        "data_obito": 14,
        "cidade": 15,
        "endereco": 16,
        "bairro": 17,
        "cep": 18,
        "uf": 19,
        "data_atualizacao": 20,
        "contatos": {
            "telefones": [21, 22, 23, 24, 25, 26, 27, 28],
            "emails": [29, 30]
        }
    }
}
```

## Funcionalidades

### 1. Processamento Paralelo
- Utiliza múltiplos processadores para processar os dados em paralelo
- O número de processadores é configurável via `num_workers` no config.json
- Cada worker processa uma parte dos dados usando LIMIT e OFFSET

### 2. Conversão Automática de Tipos
- Converte automaticamente tipos de dados do MySQL para MongoDB
- Suporta conversão de:
  - Strings (com validação UTF-8)
  - Datas (múltiplos formatos)
  - Números decimais
  - Campos opcionais
  - Arrays (telefones e emails)

### 3. Gerenciamento de Memória
- Calcula automaticamente o tamanho do lote baseado na memória disponível
- Evita sobrecarga de memória durante a migração
- Configurável via `batch_size` no config.json

### 4. Tratamento de Erros
- Logs detalhados de erros (armazenados em `tmp/logs/`)
- Tratamento de conexões perdidas
- Validação de dados durante a conversão

## Como Usar

1. Configure os arquivos `config.json` e `mapping.json` com suas credenciais e mapeamentos
2. Execute o script:
```bash
go run main.go
```

## Estrutura do Código

### internal/config
- Gerencia o carregamento e validação das configurações
- Separa configurações de conexão do mapeamento de colunas

### internal/converter
- Funções de conversão de tipos de dados
- Validação de UTF-8
- Conversão de datas e números

### internal/database
- Conexões com MySQL e MongoDB
- Gerenciamento de pools de conexão
- Tratamento de erros de conexão

### internal/migration
- Lógica de divisão de trabalho entre workers
- Cálculo de limites de memória
- Gerenciamento de lotes

### internal/models
- Estruturas de dados para documentos MongoDB
- Definição do worker de migração
- Canais de comunicação entre workers

## Requisitos

- Go 1.16 ou superior
- MySQL 5.7 ou superior
- MongoDB 4.4 ou superior

## Dependências

- `go.mongodb.org/mongo-driver/mongo` - Driver MongoDB
- `github.com/go-sql-driver/mysql` - Driver MySQL

## Segurança

- Credenciais armazenadas em arquivo de configuração
- Suporte a conexões SSL/TLS
- Validação de dados durante a conversão

## Performance

- Processamento paralelo
- Gerenciamento automático de memória
- Inserção em lotes no MongoDB
- Uso eficiente de conexões com banco de dados 