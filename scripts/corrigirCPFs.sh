#!/bin/bash

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Iniciando correção dos CPFs no MongoDB...${NC}"

# Executa o script de atualização usando mongosh
mongosh "mongodb://contatos_us:contatos_67@localhost:27017/contatos_bd" --quiet --eval '
// Função para atualizar os CPFs
async function atualizarCPFs() {
    const collection = db.pessoas;
    const total = await collection.countDocuments();
    let processados = 0;
    let atualizados = 0;
    const batchSize = 1000;
    
    print("Total de documentos a processar: " + total);
    
    // Processa em lotes para não sobrecarregar a memória
    while (processados < total) {
        const docs = await collection.find({}, { cpf: 1 }).skip(processados).limit(batchSize).toArray();
        
        for (const doc of docs) {
            if (doc.cpf) {
                const cpfStr = doc.cpf.toString();
                if (cpfStr.length < 11) {
                    const cpfCorrigido = cpfStr.padStart(11, "0");
                    await collection.updateOne(
                        { _id: doc._id },
                        { $set: { cpf: cpfCorrigido } }
                    );
                    atualizados++;
                }
            }
        }
        
        processados += docs.length;
        print("Progresso: " + processados + "/" + total + " (" + 
              Math.round((processados/total) * 100) + "%) - " + 
              atualizados + " CPFs corrigidos");
    }
    
    print("\nCorreção concluída!");
    print("Total de documentos processados: " + processados);
    print("Total de CPFs corrigidos: " + atualizados);
}

// Executa a função
atualizarCPFs();
'

# Verifica se a execução foi bem sucedida
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Script executado com sucesso!${NC}"
else
    echo -e "${RED}Erro ao executar o script${NC}"
    exit 1
fi 