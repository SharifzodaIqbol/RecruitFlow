package candidate

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Candidate struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	TelNumber string `json:"tel_number"`
	Email     string `json:"email"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
}

func GetCandidates(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query("SELECT * FROM candidate")
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	defer rows.Close()
	people := []Candidate{}
	for rows.Next() {
		person := Candidate{}
		err = rows.Scan(&person.ID, &person.Name, &person.TelNumber, &person.Email, &person.CreatedAt, &person.UpdatedAt)
		if err != nil {
			log.Println(err)
			continue
		}
		people = append(people, person)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(people)
}
func GetCandidate(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id candidate", http.StatusBadRequest)
		log.Println(err)
		return
	}
	result := db.QueryRow("SELECT * FROM candidate WHERE id = $1", id)
	var person Candidate
	err = result.Scan(&person.ID, &person.Name, &person.TelNumber, &person.Email, &person.CreatedAt, &person.UpdatedAt)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person)
}
func DeleteCandidate(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id candidate", http.StatusBadRequest)
		log.Println(err)
		return
	}
	result, err := db.Exec("DELETE FROM candidate where id = $1", id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
