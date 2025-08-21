package main

import (
	"database/sql"
	"log"

	_ "github.com/LeeDat03/gin-event-app/docs"
	"github.com/LeeDat03/gin-event-app/internal/database"
	"github.com/LeeDat03/gin-event-app/internal/env"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

//	@title			Gin Event App
//	@version		1.0
//	@description	Event management with gin framework
//	@contact.name	LeeDat03
//	@host			localhost:8000
//	@BasePath		/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter your bearer token in the format **Bearer &lt;token&gt;**

// Apply the security definition to your endpoints
//	@security	BearerAuth

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
