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
		Update[*Candidate](db, w, r)
	})
	mux.HandleFunc("DELETE /candidates/{id}", func(w http.ResponseWriter, r *http.Request) {
		Delete[*Candidate](db, w, r)
	})
	// Jobs/Postings
	mux.HandleFunc("GET /jobs", func(w http.ResponseWriter, r *http.Request) {
		Read[*Posting](db, w, r)
	})
	mux.HandleFunc("GET /jobs/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetByID[*Posting](db, w, r)
	})
	mux.HandleFunc("POST /companies/{id}/jobs", func(w http.ResponseWriter, r *http.Request) {
		Create[*Posting](db, w, r)
	})
	mux.HandleFunc("PUT /jobs/{id}", func(w http.ResponseWriter, r *http.Request) {
		Update[*Posting](db, w, r)
	})
	mux.HandleFunc("DELETE /jobs/{id}", func(w http.ResponseWriter, r *http.Request) {
		Delete[*Posting](db, w, r)
	})

	//Applications
	mux.HandleFunc("POST /applications", func(w http.ResponseWriter, r *http.Request) {
		Create[*Application](db, w, r)
	})
	mux.HandleFunc("GET /applications", func(w http.ResponseWriter, r *http.Request) {
		Read[*Application](db, w, r)
	})
	mux.HandleFunc("GET /applications/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetByID[*Application](db, w, r)
	})
	mux.HandleFunc("PUT /applications/{id}", func(w http.ResponseWriter, r *http.Request) {
		Update[*Application](db, w, r)
	})
	mux.HandleFunc("DELETE /applications/{id}", func(w http.ResponseWriter, r *http.Request) {
		Delete[*Application](db, w, r)
	})
	mux.HandleFunc("GET /jobs/{id}/applications", func(w http.ResponseWriter, r *http.Request) {
		Read[*Application](db, w, r)
	})
	mux.HandleFunc("GET /candidates/{id}/applications", func(w http.ResponseWriter, r *http.Request) {
		Read[*Application](db, w, r)
	})
	//Interviews
	mux.HandleFunc("POST /applications", func(w http.ResponseWriter, r *http.Request) {
		Create[*Application](db, w, r)
	})
	mux.HandleFunc("GET /interviews", func(w http.ResponseWriter, r *http.Request) {
		Read[*Application](db, w, r)
	})
	mux.HandleFunc("GET /interviews/{id}", func(w http.ResponseWriter, r *http.Request) {
		GetByID[*Application](db, w, r)
	})
	mux.HandleFunc("PUT /interviews/{id}", func(w http.ResponseWriter, r *http.Request) {
		Update[*Application](db, w, r)
	})
	mux.HandleFunc("DELETE /interviews/{id}", func(w http.ResponseWriter, r *http.Request) {
		Delete[*Application](db, w, r)
	})
	mux.HandleFunc("GET /jobs/{id}/interviews", func(w http.ResponseWriter, r *http.Request) {
		Read[*Application](db, w, r)
	})
	mux.HandleFunc("GET /candidates/{id}/interviews", func(w http.ResponseWriter, r *http.Request) {
		Read[*Application](db, w, r)
	})
}
