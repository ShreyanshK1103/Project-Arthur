package config

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, string, error) {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found")
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("There is no PORT in the environment")
	}
	

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	log.Println("Database connection established")

	return conn, portString, nil
}