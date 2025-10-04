package posting

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Posting struct {
	ID          int    `json:"id"`
	CompanyID   int    `json:"company_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
}

func CreateJob(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Status not allowed"})
		return
	}
	Post := Posting{}
	if err := json.NewDecoder(r.Body).Decode(&Post); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer r.Body.Close()
	result, err := db.Exec("INSERT INTO job_postings (company_id, title, description, created_at, updated_at) VALUES ($1, $2, $3, Now(), Now())",
		Post.CompanyID, Post.Title, Post.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	rowAffected, _ := result.RowsAffected()
	if rowAffected == 0 {
		http.Error(w, "Job not Added", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func GetJobs(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Status not allowed"})
		return
	}
	idStr := r.PathValue("company_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}
	Posts := []Posting{}
	rows, err := db.Query("SELECT * FROM job_postings WHERE company_id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	for rows.Next() {
		Post := Posting{}
		err := rows.Scan(&Post.ID, &Post.CompanyID, &Post.Title, &Post.Description, &Post.CreatedAt, &Post.UpdatedAt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}
		Posts = append(Posts, Post)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Posts)
}
func UpdateJob(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Status not allowed"})
		return
	}
	idStr := r.PathValue("job_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}
	var Post Posting
	if err := json.NewDecoder(r.Body).Decode(&Post); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	result, err := db.Exec("Update job_postings SET title=$2, description=$3, updated_at=Now() WHERE id = $1", id, Post.Title, Post.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	rowAffected, _ := result.RowsAffected()
	if rowAffected == 0 {
		http.Error(w, "Job not fount", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
func DeleteJob(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.PathValue("job_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid job id", http.StatusBadRequest)
		return
	}
	result, err := db.Exec("DELETE FROM job_postings WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	rowAffected, _ := result.RowsAffected()
	if rowAffected == 0 {
		http.Error(w, "Job not found", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
