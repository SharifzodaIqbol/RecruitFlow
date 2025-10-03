package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"recruitFlow/company"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
func initDB() error {
	connStr := fmt.Sprintf("user=postgres password=%s dbname=recruit sslmode=disable", os.Getenv("mypass"))
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db.Ping()
}
func main() {
	loadEnv()
	if err := initDB(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /companies", func(w http.ResponseWriter, r *http.Request) {
		company.GetCompanies(db, w, r)
	})
	mux.HandleFunc("GET /companies/{id}", func(w http.ResponseWriter, r *http.Request) {
		company.GetCompany(db, w, r)
	})
	mux.HandleFunc("POST /companies", func(w http.ResponseWriter, r *http.Request) {
		company.CreateCompanies(db, w, r)
	})
	mux.HandleFunc("PUT /companies/{id}", func(w http.ResponseWriter, r *http.Request) {
		company.UpdateCompany(db, w, r)
	})
	mux.HandleFunc("DELETE /companies/{id}", func(w http.ResponseWriter, r *http.Request) {
		company.DeleteCompany(db, w, r)
	})
	http.ListenAndServe(":8082", mux)
}
