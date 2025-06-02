package dbcommand

import (
	"database/sql"
	"fmt"
	"os"
)

func RunSeeder() error {
	viper := newViper()

	dbHost := viper.GetString("DATABASE_HOST")
	dbPort := viper.GetInt("DATABASE_PORT")
	dbUser := viper.GetString("DATABASE_USER")
	dbPass := viper.GetString("DATABASE_PASS")
	dbName := viper.GetString("DATABASE_NAME")

	connString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	seedFiles := []string{"./db/seeds/seed_rbac.sql", "./db/seeds/seed_genres.sql", "./db/seeds/seed_videos.sql", "./db/seeds/seed_user.sql"}
	var sqlString string
	for _, file := range seedFiles {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		sqlString += string(sqlBytes)
	}

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(sqlString)
	if err != nil {
		return err
	}

	return nil
}
