package helper

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Company struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
type Posting struct {
	ID          int       `json:"id"`
	CompanyID   int       `json:"company_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
}
type Candidate struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	TelNumber string    `json:"tel_number"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Reflector interface {
	GetNameDB() string
	GetParam() string
	GetPlaceholder() string
	GetValues() []interface{}
	GetFields() []interface{}
}

func (company Company) GetNameDB() string        { return "companies" }
func (company Company) GetParam() string         { return "name, created_at, updated_at" }
func (company Company) GetPlaceholder() string   { return "$1, NOW(), NOW()" }
func (company Company) GetValues() []interface{} { return []interface{}{&company.Name} }
func (company Company) GetFields() []interface{} {
	return []interface{}{
		&company.ID,
		&company.Name,
		&company.CreatedAt,
		&company.UpdatedAt,
	}
}

func (candidate Candidate) GetNameDB() string      { return "candidate" }
func (candidate Candidate) GetParam() string       { return "name,tel_number, email, created_at, updated_at" }
func (candidate Candidate) GetPlaceholder() string { return "$1, $2, $3, NOW(), NOW()" }
func (candidate Candidate) GetValues() []interface{} {
	return []interface{}{candidate.Name, candidate.TelNumber, candidate.Email}
}
func (candidate Candidate) GetFields() []interface{} {
	return []interface{}{
		&candidate.ID,
		&candidate.Name,
		&candidate.TelNumber,
		&candidate.Email,
		&candidate.CreatedAt,
		&candidate.UpdatedAt,
	}
}
func (post Posting) GetNameDB() string { return "job_postings" }
func (post Posting) GetParam() string {
	return "company_id, title, description, created_at, updated_at"
}
func (post Posting) GetPlaceholder() string { return "$1, $2, $3, NOW(), NOW()" }
func (post Posting) GetValues() []interface{} {
	return []interface{}{post.CompanyID, post.Title, post.Description}
}
func (post Posting) GetFields() []interface{} {
	return []interface{}{
		&post.ID,
		&post.CompanyID,
		&post.Title,
		&post.Description,
		&post.CreatedAt,
		&post.UpdatedAt,
	}
}
func MethodAllowed(w http.ResponseWriter, r *http.Request, MethodName string) {
	if r.Method != MethodName {
		http.Error(w, "Not allowed Method", http.StatusMethodNotAllowed)
		return
	}
}
func MethodStatus(w http.ResponseWriter, infoErr string, code int, err error) {
	http.Error(w, infoErr, code)
	log.Println(err)
}
func GetIDPath(w http.ResponseWriter, r *http.Request, nameID string) int {
	idStr := r.PathValue(nameID)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		log.Fatal()
	}
	return id
}
func Affected(w http.ResponseWriter, result sql.Result) {
	affected, _ := result.RowsAffected()
	if affected == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
}
