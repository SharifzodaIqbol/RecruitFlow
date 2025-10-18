package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"recruitFlow/pkg/store"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания мока БД: %v", err)
	}
	defer db.Close()
	company := &store.Company{Name: "Тестовая компания"}
	body, _ := json.Marshal(company)

	mock.ExpectExec("INSERT INTO companies").
		WithArgs("Тестовая компания").
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest("POST", "/companies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	store.Create[*store.Company](db, w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания выполнены: %v", err)
	}
}

func TestCreate_InvalidJSON(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	body := bytes.NewBufferString("{invalid json")

	req := httptest.NewRequest("POST", "/companies", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	store.Create[*store.Company](db, w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус 400 для невалидного JSON, получен %d", w.Code)
	}

	mock.ExpectationsWereMet()
}

func TestCreate_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	company := &store.Company{Name: "Тестовая компания"}
	body, _ := json.Marshal(company)

	mock.ExpectExec("INSERT INTO companies").
		WithArgs("Тестовая компания").
		WillReturnError(sql.ErrConnDone)

	req := httptest.NewRequest("POST", "/companies", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	store.Create[*store.Company](db, w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус 400 при ошибке БД, получен %d", w.Code)
	}
}

func TestRead_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка создания мока БД: %v", err)
	}
	defer db.Close()

	createdAt := time.Now()
	updatedAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
		AddRow(1, "Компания 1", createdAt, updatedAt).
		AddRow(2, "Компания 2", createdAt, updatedAt)

	mock.ExpectQuery("SELECT \\* FROM companies").
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/companies", nil)
	w := httptest.NewRecorder()

	store.Read[*store.Company](db, w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d. Тело ответа: %s", w.Code, w.Body.String())
	}

	var companies []store.Company
	if err := json.Unmarshal(w.Body.Bytes(), &companies); err != nil {
		t.Errorf("Ошибка разбора JSON ответа: %v", err)
	}

	if len(companies) != 2 {
		t.Errorf("Ожидалось 2 компании, получено %d", len(companies))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания выполнены: %v", err)
	}
}
