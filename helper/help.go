package helper

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func Create[T any](db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not allowed Method", http.StatusMethodNotAllowed)
		return
	}
	var item T
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer r.Body.Close()
	var result sql.Result
	var err error = nil
	switch any(item).(type) {
	case Company:
		query := `INSERT INTO companies (name, created_at, updated_at) VALUES ($1, NOW(), NOW())`
		company := any(item).(Company)
		result, err = db.Exec(query, company.Name)
	case Candidate:
		query := `INSERT INTO candidate (name,tel_number, email, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`
		candidate := any(item).(Candidate)
		result, err = db.Exec(query, candidate.Name, candidate.TelNumber, candidate.Email)
	case Posting:
		query := `INSERT INTO job_postings (company_id, title, description, created_at, updated_at) VALUES ($1, $2, $3, NOW(), NOW())`
		post := any(item).(Posting)
		result, err = db.Exec(query, post.CompanyID, post.Title, post.Description)
	}
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func Read[T any](db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Not allowed Method", http.StatusMethodNotAllowed)
		return
	}
	var item T
	var rows *sql.Rows
	var err error = nil
	var result []any
	switch any(item).(type) {
	case Company:
		company := any(item).(Company)
		rows, err = db.Query(`SELECT * FROM companies`)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		for rows.Next() {
			err = rows.Scan(&company.ID, &company.Name, &company.CreatedAt, &company.UpdatedAt)
			if err != nil {
				http.Error(w, "Server Error", http.StatusInternalServerError)
				continue
			}
			result = append(result, company)
		}
	case Candidate:
		candidate := any(item).(Candidate)
		rows, err = db.Query(`SELECT * FROM candidate`)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		for rows.Next() {
			err = rows.Scan(&candidate.ID, &candidate.Name, &candidate.TelNumber, &candidate.Email, &candidate.CreatedAt, &candidate.UpdatedAt)
			if err != nil {
				http.Error(w, "Server Error", http.StatusInternalServerError)
				continue
			}
			result = append(result, candidate)
		}
	case Posting:
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}
		post := any(item).(Posting)
		rows, err = db.Query(`SELECT * FROM job_postings WHERE company_id = $1`, id)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
		for rows.Next() {
			err = rows.Scan(&post.ID, &post.CompanyID, &post.Title, &post.Description, &post.CreatedAt, &post.UpdatedAt)
			if err != nil {
				http.Error(w, "Server Error", http.StatusInternalServerError)
				continue
			}
			result = append(result, post)
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
func GetByID[T any](db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Not allowed Method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	var item T
	switch any(item).(type) {
	case Company:
		company := any(item).(Company)
		row := db.QueryRow("SELECT * FROM companies WHERE id = $1", id)
		err = row.Scan(&company.ID, &company.Name, &company.CreatedAt, &company.UpdatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not Fount", http.StatusNotFound)
				return
			}
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(company)
	case Candidate:
		candidate := any(item).(Candidate)
		row := db.QueryRow("SELECT * FROM candidate WHERE id = $1", id)
		err = row.Scan(&candidate.ID, &candidate.Name, &candidate.TelNumber, &candidate.Email, &candidate.CreatedAt, &candidate.UpdatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not Fount", http.StatusNotFound)
				return
			}
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(candidate)
	case Posting:
		idStr = r.PathValue("job_id")
		job, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}
		post := any(item).(Posting)
		row := db.QueryRow("SELECT * FROM job_postings WHERE id = $1", job)
		err = row.Scan(&post.ID, &post.CompanyID, &post.Title, &post.Description, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not Fount", http.StatusNotFound)
				return
			}
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	}
}
func Delete[T any](db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Not allowed Method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	var item T
	var result sql.Result
	switch any(item).(type) {
	case Company:
		result, err = db.Exec(`DELETE FROM companies WHERE id = $1`, id)
		if err != nil {
			log.Println(err)
			return
		}
	case Candidate:
		result, err = db.Exec(`DELETE FROM candidate WHERE id = $1`, id)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
	case Posting:
		idStr = r.PathValue("job_id")
		job, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}
		result, err = db.Exec(`DELETE FROM job_postings WHERE id = $1`, job)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
	}
	rowAffected, _ := result.RowsAffected()
	if rowAffected == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func Update[T any](db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Not allowed Method", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	var item T
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Bad Reques", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	switch any(item).(type) {
	case Company:
		company := any(item).(Company)
		_, err = db.Exec(`UPDATE companies SET name = $1, updated_at = NOW() WHERE id = $2`, company.Name, id)
		if err != nil {
			log.Println(err)
			return
		}
	case Candidate:
		candidate := any(item).(Candidate)
		_, err = db.Exec(`UPDATE candidate SET name = $1, tel_number = $2, email = $3, updated_at = NOW() WHERE id = $4`, candidate.Name, candidate.TelNumber, candidate.Email, id)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
	case Posting:
		post := any(item).(Posting)
		_, err = db.Exec(`UPDATE job_postings SET title = $1, description = $2, updated_at = NOW() WHERE company_id = $3`, post.Title, post.Description, id)
		if err != nil {
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
