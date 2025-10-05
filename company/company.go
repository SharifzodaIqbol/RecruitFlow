package company

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Company struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
}

func GetCompanies(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Status not allowed"})
		return
	}
	rows, err := db.Query("SELECT * FROM companies")
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	Companies := []Company{}
	for rows.Next() {

		Comp := Company{}
		err := rows.Scan(&Comp.ID, &Comp.Name, &Comp.CreatedAt, &Comp.UpdatedAt)
		if err != nil {
			log.Println(err)
			continue
		}
		Companies = append(Companies, Comp)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Companies)
}

func GetCompany(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("company_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}
	row := db.QueryRow("SELECT * FROM companies where id=$1", id)
	Comp := Company{}
	err = row.Scan(&Comp.ID, &Comp.Name, &Comp.CreatedAt, &Comp.UpdatedAt)
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Comp)
}
func CreateCompanies(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusOK)
}
func UpdateCompany(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Status not allowed"})
		return
	}
	idStr := r.PathValue("company_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	var Comp Company
	if err := json.NewDecoder(r.Body).Decode(&Comp); err != nil {
		log.Fatal(err)
	}
	defer r.Body.Close()
	result, err := db.Exec("UPDATE companies SET name = $1, updated_at = Now() WHERE id = $2", Comp.Name, id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func DeleteCompany(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"err": "Status not allowed"})
		return
	}
	idStr := r.PathValue("company_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id company", http.StatusBadRequest)
		return
	}
	result, err := db.Exec("DELETE FROM companies where id = $1", id)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
