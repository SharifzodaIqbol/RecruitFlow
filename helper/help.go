package helper

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	}
	Affected(w, result)
	w.WriteHeader(http.StatusOK)
}
func Read[T Reflector](db *sql.DB, w http.ResponseWriter, r *http.Request) {
	MethodAllowed(w, r, http.MethodGet)
	var item T
	query := fmt.Sprintf("SELECT * FROM %s", item.GetNameDB())
	rows, err := db.Query(query)
	if err != nil {
		MethodStatus(w, "Server Error", http.StatusInternalServerError, err)
	}
	defer rows.Close()

	var result []interface{}

	for rows.Next() {
		var newItem T
		err = rows.Scan(newItem.GetFields()...)
		if err != nil {
			log.Printf("Scan error (Пропустили строку): %v", err)
			continue
		}
		result = append(result, newItem.GetValues()...)
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
	id := GetIDPath(w, r, "id")
	var item T
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", item.GetNameDB())
	row := db.QueryRow(query, id)
	err := row.Scan(item.GetValues()...)
	if err != nil {
		MethodStatus(w, "Not Found", http.StatusNotFound, err)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item.GetValues())
}
func Delete[T any](db *sql.DB, w http.ResponseWriter, r *http.Request) {
	MethodAllowed(w, r, http.MethodDelete)
	id := GetIDPath(w, r, "id")
	var item T
	var result sql.Result
	var err error
	switch any(item).(type) {
	case Company:
		result, err = db.Exec(`DELETE FROM companies WHERE id = $1`, id)
		if err != nil {
			MethodStatus(w, "Bad Request", http.StatusInternalServerError, err)
		}
	case Candidate:
		result, err = db.Exec(`DELETE FROM candidate WHERE id = $1`, id)
		if err != nil {
			MethodStatus(w, "Bad Request", http.StatusInternalServerError, err)
		}
	case Posting:
		id := GetIDPath(w, r, "job_id")
		result, err = db.Exec(`DELETE FROM job_postings WHERE id = $1`, id)
		if err != nil {
			MethodStatus(w, "Bad Request", http.StatusInternalServerError, err)
		}
	}
	Affected(w, result)
	w.WriteHeader(http.StatusOK)
}
func Update[T any](db *sql.DB, w http.ResponseWriter, r *http.Request) {
	MethodAllowed(w, r, http.MethodPut)
	id := GetIDPath(w, r, "id")
	var item T
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		MethodStatus(w, "Bad Request", http.StatusBadRequest, err)
	}
	defer r.Body.Close()
	switch any(item).(type) {
	case Company:
		company := any(item).(Company)
		_, err := db.Exec(`UPDATE companies SET name = $1, updated_at = NOW() WHERE id = $2`, company.Name, id)
		if err != nil {
			MethodStatus(w, "Bad Request", http.StatusBadRequest, err)
		}
	case Candidate:
		candidate := any(item).(Candidate)
		_, err := db.Exec(`UPDATE candidate SET name = $1, tel_number = $2, email = $3, updated_at = NOW() WHERE id = $4`, candidate.Name, candidate.TelNumber, candidate.Email, id)
		if err != nil {
			MethodStatus(w, "Bad Request", http.StatusBadRequest, err)
		}
	case Posting:
		post := any(item).(Posting)
		_, err := db.Exec(`UPDATE job_postings SET title = $1, description = $2, updated_at = NOW() WHERE company_id = $3`, post.Title, post.Description, id)
		if err != nil {
			MethodStatus(w, "Bad Request", http.StatusBadRequest, err)
		}
	}
	w.WriteHeader(http.StatusOK)
}
