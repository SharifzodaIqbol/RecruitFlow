package store_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"recruitFlow/pkg/store"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	godotenv.Load("../.env")

	if err := initTestDB(); err != nil {
		fmt.Printf("Failed to initialize test database: %v\n", err)
		os.Exit(1)
	}
	defer testDB.Close()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func initTestDB() error {
	password := os.Getenv("mypass")
	if password == "" {
		fmt.Println("Warning: MY_TEST_PASS not set. Using default 'mypass'.")
		password = "mypass"
	}

	connStr := fmt.Sprintf("user=postgres password=%s dbname=recruit sslmode=disable", password)

	var err error
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return testDB.Ping()
}

func clearTable(t *testing.T, tableName string) {
	_, err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tableName))
	if err != nil {
		t.Fatalf("Failed to truncate table %s: %v", tableName, err)
	}
}

func executeRequest(mux *http.ServeMux, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func TestCompanyCRUD(t *testing.T) {
	clearTable(t, "companies")

	mux := http.NewServeMux()
	store.SetupRoutes(mux, testDB)

	var companyID int

	t.Run("Create", func(t *testing.T) {
		initialName := "New Test Company"
		initialCompany := store.Company{Name: initialName}
		body, _ := json.Marshal(initialCompany)

		req := httptest.NewRequest(http.MethodPost, "/companies", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(mux, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Create failed: Expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		err := testDB.QueryRow("SELECT id FROM companies WHERE name = $1", initialName).Scan(&companyID)
		if err != nil {
			t.Fatalf("Failed to retrieve ID after creation: %v", err)
		}
	})

	if companyID == 0 {
		t.Skip("Skipping R, U, D tests because Create failed to get an ID.")
	}

	t.Run("ReadByID", func(t *testing.T) {
		url := fmt.Sprintf("/companies/%d", companyID)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = req.WithContext(context.WithValue(req.Context(), "id", fmt.Sprintf("%d", companyID)))

		rr := executeRequest(mux, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("ReadByID failed: Expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var fetchedCompany store.Company
		if err := json.NewDecoder(rr.Body).Decode(&fetchedCompany); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}
		if fetchedCompany.ID != companyID {
			t.Errorf("ID mismatch. Expected %d, got %d", companyID, fetchedCompany.ID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		updatedName := "Updated Company Name"
		updatedCompany := store.Company{Name: updatedName}
		body, _ := json.Marshal(updatedCompany)

		url := fmt.Sprintf("/companies/%d", companyID)
		req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), "id", fmt.Sprintf("%d", companyID)))

		rr := executeRequest(mux, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Update failed: Expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var name string
		err := testDB.QueryRow("SELECT name FROM companies WHERE id = $1", companyID).Scan(&name)
		if err != nil {
			t.Fatalf("Failed to verify update in DB: %v", err)
		}
		if name != updatedName {
			t.Errorf("DB name not updated. Expected %s, got %s", updatedName, name)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		url := fmt.Sprintf("/companies/%d", companyID)
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		req = req.WithContext(context.WithValue(req.Context(), "id", fmt.Sprintf("%d", companyID)))

		rr := executeRequest(mux, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("Delete failed: Expected status %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
		}

		var count int
		testDB.QueryRow("SELECT COUNT(*) FROM companies WHERE id = $1", companyID).Scan(&count)
		if count != 0 {
			t.Errorf("Company not deleted. Found %d row(s)", count)
		}
	})
}

func TestErrorHandling(t *testing.T) {
	mux := http.NewServeMux()
	store.SetupRoutes(mux, testDB)
	t.Run("InvalidID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/companies/abc", nil)
		req = req.WithContext(context.WithValue(req.Context(), "id", "abc"))
		rr := executeRequest(mux, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for invalid ID, got %d. Body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "Invalid id") {
			t.Errorf("Expected 'Invalid id' message, got: %s", rr.Body.String())
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		clearTable(t, "companies")
		nonExistentID := 9999
		url := fmt.Sprintf("/companies/%d", nonExistentID)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		req = req.WithContext(context.WithValue(req.Context(), "id", fmt.Sprintf("%d", nonExistentID)))
		rr := executeRequest(mux, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status %d for not found, got %d. Body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
		}
	})

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/companies", nil)
		rr := executeRequest(mux, req)

		if rr.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected status %d for wrong method, got %d. Body: %s", http.StatusMethodNotAllowed, rr.Code, rr.Body.String())
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		invalidBody := strings.NewReader(`{"name": "Test", "extra": }`) // Неправильный JSON
		req := httptest.NewRequest(http.MethodPost, "/companies", invalidBody)
		req.Header.Set("Content-Type", "application/json")
		rr := executeRequest(mux, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for invalid JSON, got %d. Body: %s", http.StatusBadRequest, rr.Code, rr.Body.String())
		}
		if !strings.Contains(rr.Body.String(), "Bad Request") {
			t.Errorf("Expected 'Bad Request' message, got: %s", rr.Body.String())
		}
	})
}
