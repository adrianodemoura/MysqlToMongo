#!/bin/bash
mongosh -u $1 -p $2 --eval 'use contatos_bd; db.createUser({user: "contatos_us", pwd: "contatos_67", roles: [{role: "readWrite", db: "contatos_bd"}]})' 