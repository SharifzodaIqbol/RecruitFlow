package store

import (
	"database/sql"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux, db *sql.DB) {
	// Companies
	mux.HandleFunc("GET /companies", func(w http.ResponseWriter, r *http.Request) {
		Read[*Company](db, w, r)
	})
	mux.HandleFunc("GET /companies/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetByID[*Company](db, w, r)
	})
	mux.HandleFunc("POST /companies", func(w http.ResponseWriter, r *http.Request) {
		Create[*Company](db, w, r)
	})
	mux.HandleFunc("PUT /companies/{id}", func(w http.ResponseWriter, r *http.Request) {
		Update[*Company](db, w, r)
	})
	mux.HandleFunc("DELETE /companies/{id}", func(w http.ResponseWriter, r *http.Request) {
		Delete[*Company](db, w, r)
	})

	// Jobs/Postings
	mux.HandleFunc("POST /companies/jobs", func(w http.ResponseWriter, r *http.Request) {
		Create[*Posting](db, w, r)
	})
	mux.HandleFunc("GET /companies/{id}/jobs", func(w http.ResponseWriter, r *http.Request) {
		Read[*Posting](db, w, r)
	})
	mux.HandleFunc("GET /companies/{id}/jobs/{job_id}", func(w http.ResponseWriter, r *http.Request) {
		GetByID[*Posting](db, w, r)
	})
	mux.HandleFunc("PUT /companies/{id}/jobs/{job_id}", func(w http.ResponseWriter, r *http.Request) {
		Update[*Posting](db, w, r)
	})
	mux.HandleFunc("DELETE /companies/{id}/jobs/{job_id}", func(w http.ResponseWriter, r *http.Request) {
		Delete[*Posting](db, w, r)
	})

	// Candidates
	mux.HandleFunc("POST /candidates", func(w http.ResponseWriter, r *http.Request) {
		Create[*Candidate](db, w, r)
	})
	mux.HandleFunc("GET /candidates", func(w http.ResponseWriter, r *http.Request) {
		Read[*Candidate](db, w, r)
	})
	mux.HandleFunc("GET /candidates/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetByID[*Candidate](db, w, r)
	})
	mux.HandleFunc("PUT /candidates/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetByID[*Candidate](db, w, r)
	})
	mux.HandleFunc("DELETE /candidates/{id}", func(w http.ResponseWriter, r *http.Request) {
		Delete[*Candidate](db, w, r)
	})
}
