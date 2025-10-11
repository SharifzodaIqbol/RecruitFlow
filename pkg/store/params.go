package store

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

func MethodAllowed(w http.ResponseWriter, r *http.Request, MethodName string) {
	if r.Method != MethodName {
		http.Error(w, "Not allowed Method", http.StatusMethodNotAllowed)
		return
	}
}
func MethodStatus(w http.ResponseWriter, infoErr string, code int, err error) {
	http.Error(w, infoErr, code)
	log.Fatal(err)
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
func GetIDByTypeStruct(item any, w http.ResponseWriter, r *http.Request) int {
	id := GetIDPath(w, r, "id")
	if _, ok := any(item).(*Posting); ok {
		id = GetIDPath(w, r, "job_id")
	}
	return id
}
