package main

import (
	"log"
	dbcommand "vrs-api/db/command"
)

func main() {
	if err := dbcommand.RunMigrations(); err != nil {
		log.Fatal(err)
	}
	log.Println("Migrations applied successfully")
}
