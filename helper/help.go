package helper

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type Company struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
}
type Posting struct {
	ID          int    `json:"id"`
	CompanyID   int    `json:"company_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
}
type Candidate struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	TelNumber string `json:"tel_number"`
	Email     string `json:"email"`
	CreatedAt string `json:"-"`
	UpdatedAt string `json:"-"`
}

func Create[T any](db *sql.DB, w http.ResponseWriter, r *http.Request, tableName string) {
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
	var parametrs []string
	var placeholders []string
	var values []interface{}
	counter := 0
	t := reflect.TypeOf(item)
	v := reflect.ValueOf(item)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" || jsonTag == "id" {
			continue
		}
		parametrs = append(parametrs, jsonTag)
		counter++
		placeholders = append(placeholders, fmt.Sprintf("$%d", counter))
		values = append(values, fieldValue.Interface())
	}
	parametrs = append(parametrs, "created_at", "updated_at")
	placeholders = append(placeholders, fmt.Sprintf("$%d", counter+1), fmt.Sprintf("$%d", counter+2))
	values = append(values, "NOW()", "NOW()")
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(parametrs, ", "),
		strings.Join(placeholders, ", "),
	)
	result, err := db.Exec(query, values...)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	rowAffected, _ := result.RowsAffected()
	if rowAffected == 0 {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func Read[T any](db *sql.DB, w http.ResponseWriter, r *http.Request, tableName string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Not allowed Method", http.StatusMethodNotAllowed)
		return
	}
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	results := make([]T, 0)
	for rows.Next() {
		currentItem := new(T)
		v := reflect.ValueOf(currentItem).Elem()
		numsFiled := v.NumField()
		scanArgs := make([]interface{}, numsFiled)
		for i := 0; i < numsFiled; i++ {
			scanArgs[i] = v.Field(i).Addr().Interface()
		}
		if err := rows.Scan(scanArgs...); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, *currentItem)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}
