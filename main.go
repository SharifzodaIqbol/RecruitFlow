package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"recruitFlow/helper"

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
		helper.Read[helper.Company](db, w, r)
	})
	mux.HandleFunc("GET /companies/{id}", func(w http.ResponseWriter, r *http.Request) {
		helper.GetByID[helper.Company](db, w, r)
	})
	mux.HandleFunc("POST /companies", func(w http.ResponseWriter, r *http.Request) {
		helper.Create[helper.Company](db, w, r)
	})
	mux.HandleFunc("PUT /companies/{id}", func(w http.ResponseWriter, r *http.Request) {
		helper.Update[helper.Company](db, w, r)
	})
	mux.HandleFunc("DELETE /companies/{id}", func(w http.ResponseWriter, r *http.Request) {
		helper.Delete[helper.Company](db, w, r)
	})
	mux.HandleFunc("POST /companies/jobs", func(w http.ResponseWriter, r *http.Request) {
		helper.Create[helper.Posting](db, w, r)
	})
	mux.HandleFunc("GET /companies/{id}/jobs", func(w http.ResponseWriter, r *http.Request) {
		helper.Read[helper.Posting](db, w, r)
	})
	mux.HandleFunc("GET /companies/{id}/jobs/{job_id}", func(w http.ResponseWriter, r *http.Request) {
		helper.GetByID[helper.Posting](db, w, r)
	})
	mux.HandleFunc("PUT /companies/{id}/jobs/{job_id}", func(w http.ResponseWriter, r *http.Request) {
		helper.Update[helper.Posting](db, w, r)
	})
	mux.HandleFunc("DELETE /companies/{id}/jobs/{job_id}", func(w http.ResponseWriter, r *http.Request) {
		helper.Delete[helper.Posting](db, w, r)
	})
	mux.HandleFunc("POST /candidates", func(w http.ResponseWriter, r *http.Request) {
		helper.Create[helper.Candidate](db, w, r)
	})
	mux.HandleFunc("GET /candidates", func(w http.ResponseWriter, r *http.Request) {
		helper.Read[helper.Candidate](db, w, r)
	})
	mux.HandleFunc("GET /candidates/{id}", func(w http.ResponseWriter, r *http.Request) {
		helper.GetByID[helper.Candidate](db, w, r)
	})
	mux.HandleFunc("DELETE /candidates/{id}", func(w http.ResponseWriter, r *http.Request) {
		helper.Delete[helper.Candidate](db, w, r)
	})
	http.ListenAndServe(":8082", mux)
}
