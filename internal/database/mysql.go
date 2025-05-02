package database

import (
	"database/sql"
	"fmt"

	"MysqlToMongo/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

func ConnectMySQL(config *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.MySQL.User,
		config.MySQL.Password,
		config.MySQL.Host,
		config.MySQL.Port,
		config.MySQL.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao MySQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao testar conexão com MySQL: %v", err)
	}

	return db, nil
}
