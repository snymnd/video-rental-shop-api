package main

import (
	"flag"
	"log"
	dbcommand "vrs-api/db/command"
)

func main() {
	var isDown bool
	flag.BoolVar(&isDown, "down", false, "Run migration down")
	flag.Parse()

	if err := dbcommand.RunMigrations(isDown); err != nil {
		log.Fatal(err)
	}
	log.Println("migrations applied successfully")

	if !isDown {
		if err := dbcommand.RunSeeder(); err != nil {
			log.Fatal(err)
		}
		log.Println("seeds applied successfully")
	}
}
