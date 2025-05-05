#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para calcular o tempo decorrido
tempo_decorrido() {
    local inicio=$1
    local fim=$2
    # Calcula a diferença em segundos
    local tempo_total=$(echo "$fim - $inicio" | bc)
    
    # Calcula horas, minutos, segundos e milissegundos
    local horas=$(echo "scale=0; $tempo_total / 3600" | bc)
    local minutos=$(echo "scale=0; ($tempo_total % 3600) / 60" | bc)
    local segundos=$(echo "scale=0; $tempo_total % 60" | bc)
    local milissegundos=$(echo "scale=0; ($tempo_total * 1000) % 1000" | bc)
    
    # Converte para inteiros para o printf
    horas=${horas%.*}
    minutos=${minutos%.*}
    segundos=${segundos%.*}
    milissegundos=${milissegundos%.*}
    
    # Formata a saída
    printf "Tempo de execução: %02d:%02d:%02d.%03d\n" "$horas" "$minutos" "$segundos" "$milissegundos"
}

# Verifica se mongosh está instalado
if ! command -v mongosh &> /dev/null; then
    echo -e "${RED}Erro: mongosh não está instalado. Por favor, instale o MongoDB Shell.${NC}"
    exit 1
fi

# Verifica se os parâmetros foram fornecidos
if [ $# -lt 2 ] || [ $# -gt 2 ]; then
    echo -e "${YELLOW}Uso: $0 <campo> <valor>${NC}"
    echo -e "Exemplos:"
    echo -e "  $0 telefone 31996320718"
    echo -e "  $0 cpf 12345678900"
    echo -e "  $0 email joao@email.com"
    echo -e "  $0 nome \"João Silva\""
    exit 1
fi

CAMPO=$1
VALOR=$2

# Define o campo de busca baseado no parâmetro
case $CAMPO in
    "telefone")
        CAMPO_BUSCA="contatos.telefones"
        OPERADOR="\$in"
        ;;
    "email")
        CAMPO_BUSCA="contatos.emails"
        OPERADOR="\$in"
        ;;
    "cpf")
        CAMPO_BUSCA="cpf"
        OPERADOR="="
        ;;
    "nome")
        CAMPO_BUSCA="nome"
        OPERADOR="regex"
        ;;
    *)
        echo -e "${RED}Campo inválido. Campos disponíveis: telefone, email, cpf, nome${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}Buscando por $CAMPO: $VALOR${NC}"

# Inicia o cronômetro com precisão de milissegundos
inicio=$(date +%s.%N)

# Executa a busca usando mongosh
case $OPERADOR in
    "\$in")
        # Busca em arrays usando $in
        mongosh "mongodb://contatos_us:contatos_67@localhost:27017/contatos_bd" --quiet --eval "
        db.pessoas.find(
            { '$CAMPO_BUSCA': { \$in: ['$VALOR'] } },
            {
                cpf: 1,
                nome: 1,
                'contatos.telefones': 1,
                'contatos.emails': 1,
                data_atualizacao: {
                    \$dateToString: {
                        date: '\$data_atualizacao',
                        format: '%d/%m/%Y %H:%M:%S',
                        timezone: 'America/Sao_Paulo'
                    }
                },
                _id: 0
            }
        ).pretty()
        "
        ;;
    "regex")
        # Busca por nome usando regex case-insensitive
        mongosh "mongodb://contatos_us:contatos_67@localhost:27017/contatos_bd" --quiet --eval "
        db.pessoas.find(
            { '$CAMPO_BUSCA': { \$regex: '^$VALOR$', \$options: 'i' } },
            {
                cpf: 1,
                nome: 1,
                'contatos.telefones': 1,
                'contatos.emails': 1,
                data_atualizacao: {
                    \$dateToString: {
                        date: '\$data_atualizacao',
                        format: '%d/%m/%Y %H:%M:%S',
                        timezone: 'America/Sao_Paulo'
                    }
                },
                _id: 0
            }
        ).pretty()
        "
        ;;
    *)
        # Busca direta para campos não-array
        mongosh "mongodb://contatos_us:contatos_67@localhost:27017/contatos_bd" --quiet --eval "
        db.pessoas.find(
            { '$CAMPO_BUSCA': '$VALOR' },
            {
                cpf: 1,
                nome: 1,
                'contatos.telefones': 1,
                'contatos.emails': 1,
                data_atualizacao: {
                    \$dateToString: {
                        date: '\$data_atualizacao',
                        format: '%d/%m/%Y %H:%M:%S',
                        timezone: 'America/Sao_Paulo'
                    }
                },
                _id: 0
            }
        ).pretty()
        "
        ;;
esac

# Finaliza o cronômetro e calcula o tempo
fim=$(date +%s.%N)
tempo_decorrido $inicio $fim

# Verifica se a busca retornou resultados
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Busca concluída${NC}"
else
    echo -e "${RED}Erro ao executar a busca${NC}"
    exit 1
fi 