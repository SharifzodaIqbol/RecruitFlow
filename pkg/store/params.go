package store

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

func MethodAllowed(w http.ResponseWriter, r *http.Request, MethodName string) {
	if r.Method != MethodName {
		http.Error(w, "Not allowed Method", http.StatusMethodNotAllowed)
		return
	}
}
func GetIDPath(w http.ResponseWriter, r *http.Request, nameID string) (int, bool) {
	idStr := r.PathValue(nameID)
	if idStr == "" {
		return 0, false
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}
func Affected(w http.ResponseWriter, result sql.Result) error {
	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	if affected == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return fmt.Errorf("no rows affected")
	}
	return nil
}
