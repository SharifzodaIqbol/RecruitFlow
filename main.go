package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Company struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
}

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
func GetCompanies(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	rows, err := db.Query("SELECT * FROM companies")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	Companies := []Company{}
	for rows.Next() {

		Comp := Company{}
		err := rows.Scan(&Comp.ID, &Comp.Name, &Comp.CreatedAt, &Comp.UpdatedAt)
		if err != nil {
			fmt.Println(err)
			continue
		}
		Companies = append(Companies, Comp)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Companies)
}

func GetCompany(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}
	row := db.QueryRow("SELECT * FROM companies where id=$1", id)
	Comp := Company{}
	err = row.Scan(&Comp.ID, &Comp.Name, &Comp.CreatedAt, &Comp.UpdatedAt)
	if err != nil {
		fmt.Println(err)
		return
	}
	json.NewEncoder(w).Encode(Comp)
}
func CreateCompanies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Status not allowed"})
		return
	}
	var Comp Company
	if err := json.NewDecoder(r.Body).Decode(&Comp); err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	_, err := db.Exec("INSERT INTO companies (name, created_at, updated_at) Values ($1, Now(), Now())", Comp.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(Comp)
}
func main() {
	loadEnv()
	if err := initDB(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /companies", GetCompanies)
	mux.HandleFunc("GET /companies/{id}", GetCompany)
	mux.HandleFunc("POST /companies", CreateCompanies)
	http.ListenAndServe(":8082", mux)
}
