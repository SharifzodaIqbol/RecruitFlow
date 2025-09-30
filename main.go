package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Company struct {
	id   int
	name string
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
func information(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Все норм")
}
func companies(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	connStr := fmt.Sprintf("user=postgres password=%s dbname=recruit sslmode=disable", os.Getenv("mypass"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("SELECT * FROM companies")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	Companies := []Company{}
	for rows.Next() {
		Comp := Company{}
		err := rows.Scan(&Comp.id, &Comp.name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		Companies = append(Companies, Comp)
	}
	for _, c := range Companies {
		fmt.Fprintln(w, c.id, c.name)
	}
}
func main() {
	loadEnv()
	http.HandleFunc("/", information)
	http.HandleFunc("/companies", companies)
	http.ListenAndServe(":8082", nil)
}
