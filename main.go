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
func information(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Все норм")
}
func GetCompanies(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	rows, err := db.Query("SELECT * FROM companies")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	if !rows.Next() {
		w.Write([]byte("{}"))
		return
	}
	for rows.Next() {

		Comp := Company{}
		err := rows.Scan(&Comp.ID, &Comp.Name, &Comp.CreatedAt, &Comp.UpdatedAt)
		if err != nil {
			fmt.Println(err)
			continue
		}
		jsonData, err := json.Marshal(Comp)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintln(w, string(jsonData))
	}
}

func GetCompany(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err)
	}
	row := db.QueryRow("SELECT * FROM companies where id=$1", id)
	Comp := Company{}
	err = row.Scan(&Comp.ID, &Comp.Name, &Comp.CreatedAt, &Comp.UpdatedAt)
	if err == sql.ErrNoRows {
		w.Write([]byte("{}"))
		return
	}
	jsonData, err := json.Marshal(Comp)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(jsonData)
}
func main() {
	loadEnv()
	if err := initDB(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	http.HandleFunc("/", information)
	http.HandleFunc("/companies", GetCompanies)
	http.HandleFunc("/companies/{id}", GetCompany)
	http.ListenAndServe(":8082", nil)
}
