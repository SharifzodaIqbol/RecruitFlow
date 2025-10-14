package tests

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"recruitFlow/pkg/store"
	"testing"
)

func TestMethodAllowed(t *testing.T) {
	tests := []struct {
		name           string
		requestMethod  string
		allowedMethod  string
		expectedStatus int
	}{
		{"GET allowed", "GET", "GET", http.StatusOK},
		{"POST rejected", "POST", "GET", http.StatusMethodNotAllowed},
		{"PUT rejected", "PUT", "GET", http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.requestMethod, "/", nil)
			w := httptest.NewRecorder()

			store.MethodAllowed(w, r, tt.allowedMethod)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
func setupRequest(idValue string) *http.Request {
	r := httptest.NewRequest("GET", "/", nil)
	if idValue != "" {
		r.SetPathValue("id", idValue)
	}
	return r
}
func TestGetIDPath(t *testing.T) {
	tests := []struct {
		name       string
		idValue    string
		wantID     int
		wantExist  bool
		wantStatus int
	}{
		{
			name:       "valid id",
			idValue:    "14",
			wantID:     14,
			wantExist:  true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "id missing from path",
			idValue:    "",
			wantID:     0,
			wantExist:  false,
			wantStatus: http.StatusOK,
		},
		{
			name:       "another valid id",
			idValue:    "9999",
			wantID:     9999,
			wantExist:  true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "invalid id - not a number",
			idValue:    "abc",
			wantID:     0,
			wantExist:  false,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "negative id",
			idValue:    "-5",
			wantID:     -5,
			wantExist:  true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "empty string id",
			idValue:    "",
			wantID:     0,
			wantExist:  false,
			wantStatus: http.StatusOK,
		},
		{
			name:       "float number as id",
			idValue:    "3.14",
			wantID:     0,
			wantExist:  false,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupRequest(tt.idValue)
			w := httptest.NewRecorder()

			gotID, gotExist := store.GetIDPath(w, r, "id")
			if gotID != tt.wantID {
				t.Errorf("GetIDPath() gotID = %v, want %v", gotID, tt.wantID)
			}
			if gotExist != tt.wantExist {
				t.Errorf("GetIDPath() gotExist = %v, want %v", gotExist, tt.wantExist)
			}

			if tt.wantStatus != http.StatusOK {
				if w.Code != tt.wantStatus {
					t.Errorf("GetIDPath() status code = %v, want %v", w.Code, tt.wantStatus)
				}
			}
		})
	}
}

type mockResult struct {
	affected int64
	err      error
}

func (m mockResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (m mockResult) RowsAffected() (int64, error) {
	return m.affected, m.err
}
func TestAffected_Original(t *testing.T) {
	tests := []struct {
		name       string
		mockResult sql.Result
		wantStatus int
	}{
		{
			name:       "rows affected - original returns 200",
			mockResult: mockResult{affected: 1},
			wantStatus: 200,
		},
		{
			name:       "no rows affected - original returns 404",
			mockResult: mockResult{affected: 0},
			wantStatus: 404,
		},
		{
			name:       "error case - original has bug",
			mockResult: mockResult{err: errors.New("db error")},
			wantStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			err := store.Affected(w, tt.mockResult)
			t.Logf("Error: %v", err)

			if w.Code != tt.wantStatus {
				t.Errorf("Status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
