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
    # Usa bc para calcular a diferença em milissegundos
    local tempo_ms=$(echo "scale=3; ($fim - $inicio) * 1000" | bc)
    echo "Tempo de execução: ${tempo_ms} milissegundos"
}

# Verifica se mongosh está instalado
if ! command -v mongosh &> /dev/null; then
    echo -e "${RED}Erro: mongosh não está instalado. Por favor, instale o MongoDB Shell.${NC}"
    exit 1
fi

# Verifica se os parâmetros foram fornecidos
if [ $# -lt 2 ] || [ $# -gt 3 ]; then
    echo -e "${YELLOW}Uso: $0 <campo> <valor> [--case-sensitive]${NC}"
    echo -e "Exemplos:"
    echo -e "  $0 telefone 31996320718"
    echo -e "  $0 cpf 12345678900"
    echo -e "  $0 email joao@email.com"
    echo -e "  $0 nome \"João Silva\""
    echo -e "  $0 nome \"João Silva\" --case-sensitive"
    exit 1
fi

CAMPO=$1
VALOR=$2
CASE_SENSITIVE=false

# Verifica se a opção case-sensitive foi fornecida
if [ $# -eq 3 ] && [ "$3" = "--case-sensitive" ]; then
    CASE_SENSITIVE=true
fi

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
        OPERADOR="="
        ;;
    *)
        echo -e "${RED}Campo inválido. Campos disponíveis: telefone, email, cpf, nome${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}Buscando por $CAMPO: $VALOR${NC}"
if [ "$CASE_SENSITIVE" = true ]; then
    echo -e "${YELLOW}Modo case-sensitive ativado${NC}"
fi

# Inicia o cronômetro com precisão de milissegundos
inicio=$(date +%s.%N)

# Executa a busca usando mongosh
if [ "$OPERADOR" = "\$in" ]; then
    # Busca em arrays usando $in
    if [ "$CASE_SENSITIVE" = true ]; then
        echo -e "${YELLOW}Executando busca case-sensitive em array...${NC}"
        mongosh "mongodb://contatos_us:contatos_67@localhost:27017/contatos_bd" --quiet --eval "
        db.pessoas.find(
            { '$CAMPO_BUSCA': { \$elemMatch: { \$eq: '$VALOR' } } },
            {
                cpf: 1,
                nome: 1,
                'contatos.telefones': 1,
                'contatos.emails': 1,
                data_atualizacao: 1,
                _id: 0
            }
        ).pretty()
        "
    else
        mongosh "mongodb://contatos_us:contatos_67@localhost:27017/contatos_bd" --quiet --eval "
        db.pessoas.find(
            { '$CAMPO_BUSCA': { \$in: ['$VALOR'] } },
            {
                cpf: 1,
                nome: 1,
                'contatos.telefones': 1,
                'contatos.emails': 1,
                data_atualizacao: 1,
                _id: 0
            }
        ).pretty()
        "
    fi
else
    # Busca direta para campos não-array
    if [ "$CASE_SENSITIVE" = true ]; then
        echo -e "${YELLOW}Executando busca case-sensitive...${NC}"
        mongosh "mongodb://contatos_us:contatos_67@localhost:27017/contatos_bd" --quiet --eval "
        db.pessoas.find(
            { '$CAMPO_BUSCA': { \$eq: '$VALOR' } },
            {
                cpf: 1,
                nome: 1,
                'contatos.telefones': 1,
                'contatos.emails': 1,
                data_atualizacao: 1,
                _id: 0
            }
        ).pretty()
        "
    else
        mongosh "mongodb://contatos_us:contatos_67@localhost:27017/contatos_bd" --quiet --eval "
        db.pessoas.find(
            { '$CAMPO_BUSCA': '$VALOR' },
            {
                cpf: 1,
                nome: 1,
                'contatos.telefones': 1,
                'contatos.emails': 1,
                data_atualizacao: 1,
                _id: 0
            }
        ).pretty()
        "
    fi
fi

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