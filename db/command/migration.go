package dbcommand

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies database migrations to the specified database
func RunMigrations() error {
	viper := newViper()

	dbHost := viper.GetString("DATABASE_HOST")
	dbPort := viper.GetInt("DATABASE_PORT")
	dbUser := viper.GetString("DATABASE_USER")
	dbPass := viper.GetString("DATABASE_PASS")
	dbName := viper.GetString("DATABASE_NAME")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	fmt.Println(dbURL, "dburl")
	m, err := migrate.New("file://"+"./db/migrations", dbURL)
	if err != nil {
		fmt.Println("1")
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	srcErr, dbErr := m.Close()
	if srcErr != nil {
		fmt.Println("3")
		return srcErr
	}
	if dbErr != nil {
		fmt.Println("4")
		return dbErr
	}

	return nil
}
