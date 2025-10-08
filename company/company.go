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

func UpdateCompany(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Status not allowed"})
		return
	}
	idStr := r.PathValue("id")
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
