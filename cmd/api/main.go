package main

import (
	"database/sql"
	"log"

	"github.com/LeeDat03/gin-event-app/internal/database"
	"github.com/LeeDat03/gin-event-app/internal/env"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	models := database.NewModels(db)

	app := &application{
		port:      env.GetEnvInt("PORT", 8000),
		jwtSecret: env.GetEnvString("JWT_SECRET", "secret-123123"),
		models:    models,
	}

	if err := serve(app); err != nil {
		log.Fatal(err)
	}

}
