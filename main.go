package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var database *sql.DB

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load .env file: ", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not configured")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Database initialization failed: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Database connection failed: ", err)
	}

	database = db
	log.Println("Connected to PostgreSQL")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, "/products", http.StatusSeeOther)
	})

	mux.HandleFunc("/products", ProductsHandler)

	log.Println("Server started at http://localhost:8181")
	log.Fatal(http.ListenAndServe("127.0.0.1:8181", mux))
}
