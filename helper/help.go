package helper

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
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
	columns, err := rows.Columns()
	if err != nil {
		http.Error(w, "Error getting columns", http.StatusInternalServerError)
	}
	var results []T
	t := reflect.TypeOf((*T)(nil)).Elem()
	fieldMap := make(map[string]int)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			fieldMap[jsonTag] = i
		}
	}
	fmt.Println(fieldMap)
	for rows.Next() {
		scanArgs := make([]interface{}, len(columns))
		var item T
		v := reflect.ValueOf(&item).Elem()
		for i, col := range columns {
			if fieldIndex, exists := fieldMap[col]; exists {
				scanArgs[i] = v.Field(fieldIndex).Addr().Interface()
			} else {
				var dummy interface{}
				scanArgs[i] = &dummy
			}
		}
		if err := rows.Scan(scanArgs...); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}

		results = append(results, item)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}
func GetByID[T any](db *sql.DB, w http.ResponseWriter, r *http.Request, tableName string) {
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
	query := fmt.Sprintf("SELECT * from %s WHERE id = %d", tableName, id)
	row := db.QueryRow(query)
	var item T
	t := reflect.TypeOf((*T)(nil)).Elem()
	n := t.NumField()
	fieldsMap := make(map[string]int)
	for i := 0; i < n; i++ {
		jsonTag := t.Field(i).Tag.Get("json")
		if jsonTag != "" || jsonTag != "-" {
			fieldsMap[jsonTag] = i
		}
	}
	scanArgs := make([]interface{}, n)
	v := reflect.ValueOf(&item).Elem()
	for i := 0; i < n; i++ {
		jsonTag := t.Field(i).Tag.Get("json")
		if filedIndex, exist := fieldsMap[jsonTag]; exist {
			scanArgs[i] = v.Field(filedIndex).Addr().Interface()
		} else {
			var dummy interface{}
			scanArgs[i] = &dummy
		}
	}
	if err := row.Scan(scanArgs...); err != nil {
		http.Error(w, "Invalid id: "+err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
}
func Delete[T any](db *sql.DB, w http.ResponseWriter, r *http.Request, tableName string) {
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
	query := fmt.Sprintf("DELETE FROM %s WHERE id = %d", tableName, id)
	result, err := db.Exec(query)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	rowAffected, _ := result.RowsAffected()
	if rowAffected == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func Update[T any](db *sql.DB, w http.ResponseWriter, r *http.Request, tableName string) {
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
	t := reflect.TypeOf(item)
	v := reflect.ValueOf(item)
	n := t.NumField()
	fieldsName := []string{}
	var values []interface{}
	counter := 0
	for i := 0; i < n; i++ {
		jsonTag := t.Field(i).Tag.Get("json")
		if jsonTag == "-" {
			continue
		}
		fieldsName = append(fieldsName, jsonTag)
		values = append(values, v.Field(i).Interface())
	}
	query := fmt.Sprintf("UPDATE %s set %s=%s WHERE id = %d")
	result, err := db.Exec()
}
