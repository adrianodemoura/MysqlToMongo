#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Diretório de destino
DESTINO="/media/adrianoc/Dri235Gb/MysqlToMongo.zip"

echo -e "${GREEN}Iniciando exportação da branch MAIN...${NC}"

# Verifica se o diretório de destino existe
if [ ! -d "/media/adrianoc/Dri235Gb" ]; then
    echo -e "${RED}Erro: Diretório de destino não encontrado: /media/adrianoc/Dri235Gb${NC}"
    exit 1
fi

# Verifica se o arquivo ZIP já existe e remove
if [ -f "$DESTINO" ]; then
    echo -e "${YELLOW}Arquivo ZIP já existe. Removendo...${NC}"
    rm "$DESTINO"
fi

# Cria o arquivo ZIP excluindo arquivos desnecessários
echo -e "${GREEN}Criando arquivo ZIP...${NC}"
git archive --format=zip --output="$DESTINO" main

# Verifica se o arquivo foi criado com sucesso
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Exportação concluída com sucesso!${NC}"
    echo -e "Arquivo salvo em: ${YELLOW}$DESTINO${NC}"
    
    # Mostra o tamanho do arquivo
    TAMANHO=$(du -h "$DESTINO" | cut -f1)
    echo -e "Tamanho do arquivo: ${YELLOW}$TAMANHO${NC}"
else
    echo -e "${RED}Erro ao criar o arquivo ZIP${NC}"
    exit 1
fi 