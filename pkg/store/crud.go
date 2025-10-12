// Package store содержит общие утилиты для работы с БД и HTTP.
package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func Create[T Reflector](db *sql.DB, w http.ResponseWriter, r *http.Request) {
	MethodAllowed(w, r, http.MethodPost)
	var item T
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		MethodStatus(w, "Bad Request", http.StatusBadRequest, err)
	}
	defer r.Body.Close()
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		item.GetNameDB(),
		item.GetParam(),
		item.GetPlaceholder())
	result, err := db.Exec(query, item.GetValues()...)
	if err != nil {
		MethodStatus(w, "Bad Request", http.StatusBadRequest, err)
		return
	}
	Affected(w, result)
	w.WriteHeader(http.StatusOK)
}
func Read[T Reflector](db *sql.DB, w http.ResponseWriter, r *http.Request) {
	MethodAllowed(w, r, http.MethodGet)
	var item T
	query := fmt.Sprintf("SELECT * FROM %s", item.GetNameDB())
	if id, ok := GetIDPath(w, r, "id"); ok {
		query = fmt.Sprintf("SELECT * FROM %s WHERE id = %d", item.GetNameDB(), id)
	}
	rows, err := db.Query(query)
	if err != nil {
		MethodStatus(w, "Server Error", http.StatusInternalServerError, err)
		return
	}
	defer rows.Close()

	var result []Reflector

	for rows.Next() {
		newItem := item.New()
		err := rows.Scan(newItem.GetFields()...)
		if err != nil {
			log.Printf("Scan error (Пропустили строку): %v", err)
			continue
		}
		result = append(result, newItem)
	}
	if err = rows.Err(); err != nil {
		MethodStatus(w, "Server Error", http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
func GetByID[T Reflector](db *sql.DB, w http.ResponseWriter, r *http.Request) {
	MethodAllowed(w, r, http.MethodGet)
	var item T
	id, _ := GetIDPath(w, r, "id")
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", item.GetNameDB())
	row := db.QueryRow(query, id)
	newItem := item.New()
	err := row.Scan(newItem.GetFields()...)
	if err != nil {
		MethodStatus(w, "Not Found", http.StatusNotFound, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newItem)
}
func Delete[T Reflector](db *sql.DB, w http.ResponseWriter, r *http.Request) {
	MethodAllowed(w, r, http.MethodDelete)
	var item T
	id, _ := GetIDPath(w, r, "id")
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", item.GetNameDB())
	result, err := db.Exec(query, id)
	if err != nil {
		MethodStatus(w, "Bad Request", http.StatusBadRequest, err)
		return
	}
	Affected(w, result)
	w.WriteHeader(http.StatusOK)
}
func Update[T Reflector](db *sql.DB, w http.ResponseWriter, r *http.Request) {
	MethodAllowed(w, r, http.MethodPut)
	var item T
	id, _ := GetIDPath(w, r, "id")
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		MethodStatus(w, "Bad Request", http.StatusBadRequest, err)
	}
	defer r.Body.Close()
	param := strings.Split(item.GetParam(), ", ")
	placeholder := strings.Split(item.GetPlaceholder(), ", ")
	setParam := ""
	n := len(param) - 1
	for i := 0; i < n; i++ {
		if param[i] == "created_at" || param[i] == "id" {
			continue
		}
		setParam += param[i] + " = " + placeholder[i] + ", "
	}
	setParam += param[n] + " = " + placeholder[n]
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = %d", item.GetNameDB(), setParam, id)
	fmt.Println(query)
	result, err := db.Exec(query, item.GetValues()...)
	if err != nil {
		MethodStatus(w, "Bad Request", http.StatusBadRequest, err)
		return
	}
	Affected(w, result)
	w.WriteHeader(http.StatusOK)
}
