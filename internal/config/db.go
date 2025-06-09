package config

import (
	"database/sql"
	"fmt"
	"vrs-api/internal/util/logger"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
)

func NewDbConnection(viper *viper.Viper) *sql.DB {
	log := logger.GetLogger()
	dbHost := viper.GetString("DATABASE_HOST")
	dbPort := viper.GetInt("DATABASE_PORT")
	dbUser := viper.GetString("DATABASE_USER")
	dbPass := viper.GetString("DATABASE_PASS")
	dbName := viper.GetString("DATABASE_NAME")

	connString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sql.Open("pgx", connString)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping to databse: %v", err)
	}

	return db
}

func CloseDB(dbConn *sql.DB) {
	err := dbConn.Close()
	log := logger.GetLogger()

	if err != nil {
		log.Fatalf("got error when closing the DB connection", err)
	}
}
