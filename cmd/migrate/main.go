package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		migrationsPath = flag.String("path", "./migrations", "Path to migrations")
		databaseURL    = flag.String("database", "", "Database URL")
		action         = flag.String("action", "", "Action to perform: up, down, reset, version")
		steps          = flag.Int("steps", 0, "Number of steps for up/down actions")
	)

	flag.Parse()

	if *databaseURL == "" {
		log.Fatal("Database URL is required")
	}

	m, err := migrate.New(
		"file://"+*migrationsPath,
		*databaseURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	switch *action {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
	case "reset":
		err = m.Drop()
		if err != nil {
			log.Printf("Drop failed: %v", err)
		}
		err = m.Up()
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		fmt.Printf("Version: %d, Dirty: %t\n", version, dirty)
		return
	default:
		log.Fatal("Invalid action. Use: up, down, reset, version")
	}

	if err != nil {
		if err == migrate.ErrNoChange {
			fmt.Println("No changes")
		} else {
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		fmt.Println("Migration completed successfully")
	}
}
