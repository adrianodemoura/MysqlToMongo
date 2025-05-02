#!/bin/bash
# Conecta como root para ter permissão de excluir o banco
mongosh 'mongodb://root:Mongo6701!@localhost:27017/admin' --eval '
try {
    db = db.getSiblingDB("contatos_bd");
    db.dropDatabase();
    print("Banco contatos_bd removido com sucesso!");
} catch (error) {
    print("Erro ao remover banco: " + error);
}' 