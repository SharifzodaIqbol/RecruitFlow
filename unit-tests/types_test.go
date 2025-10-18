package tests

import (
	"recruitFlow/pkg/store"
	"testing"
	"time"
)

func TestCompanyReflectorMethods(t *testing.T) {
	company := &store.Company{
		ID:   1,
		Name: "Тестовая компания",
	}

	if company.GetNameDB() != "companies" {
		t.Errorf("Ожидалось 'companies', получено '%s'", company.GetNameDB())
	}

	expectedParam := "name, created_at, updated_at"
	if company.GetParam() != expectedParam {
		t.Errorf("Ожидалось '%s', получено '%s'", expectedParam, company.GetParam())
	}

	if company.GetPlaceholder() != "$1, NOW(), NOW()" {
		t.Errorf("Неожиданный плейсхолдер: %s", company.GetPlaceholder())
	}

	values := company.GetValues()
	if len(values) != 1 {
		t.Errorf("Ожидалось 1 значение, получено %d", len(values))
	}
	if values[0] != "Тестовая компания" {
		t.Errorf("Ожидалось 'Тестовая компания', получено '%v'", values[0])
	}

	fields := company.GetFields()
	if len(fields) != 4 {
		t.Errorf("Ожидалось 4 поля, получено %d", len(fields))
	}

	newCompany := company.New()
	if _, ok := newCompany.(*store.Company); !ok {
		t.Errorf("New() должен возвращать *Company")
	}
}

func TestCandidateReflectorMethods(t *testing.T) {
	candidate := &store.Candidate{
		Name:      "Иван Иванов",
		TelNumber: "123456789",
		Email:     "ivan@example.com",
	}

	if candidate.GetNameDB() != "candidates" {
		t.Errorf("Неверное имя таблицы: %s", candidate.GetNameDB())
	}

	expectedParam := "name, tel_number, email, created_at, updated_at"
	if candidate.GetParam() != expectedParam {
		t.Errorf("Неверные параметры: %s", candidate.GetParam())
	}

	values := candidate.GetValues()
	if len(values) != 3 {
		t.Errorf("Ожидалось 3 значения, получено %d", len(values))
	}
	if values[0] != "Иван Иванов" {
		t.Errorf("Неверное имя кандидата: %v", values[0])
	}
	if values[1] != "123456789" {
		t.Errorf("Неверный номер телефона: %v", values[1])
	}
	if values[2] != "ivan@example.com" {
		t.Errorf("Неверный email: %v", values[2])
	}
}

func TestPostingReflectorMethods(t *testing.T) {
	posting := &store.Posting{
		CompanyID:   1,
		Title:       "Инженер-программист",
		Description: "Описание вакансии",
	}

	if posting.GetNameDB() != "job_postings" {
		t.Errorf("Неверное имя таблицы для вакансии: %s", posting.GetNameDB())
	}

	values := posting.GetValues()
	if values[0] != 1 {
		t.Errorf("Неверный CompanyID: ожидалось 1, получено %v", values[0])
	}
	if values[1] != "Инженер-программист" {
		t.Errorf("Неверный заголовок: %v", values[1])
	}
	if values[2] != "Описание вакансии" {
		t.Errorf("Неверное описание: %v", values[2])
	}
}

func TestApplicationReflectorMethods(t *testing.T) {
	app := &store.Application{
		JobID:       1,
		CandidateID: 2,
		Status:      "ожидание",
	}

	if app.GetNameDB() != "applications" {
		t.Errorf("Неверное имя таблицы для заявки: %s", app.GetNameDB())
	}

	values := app.GetValues()
	if values[0] != 1 {
		t.Errorf("Неверный JobID: ожидалось 1, получено %v", values[0])
	}
	if values[1] != 2 {
		t.Errorf("Неверный CandidateID: ожидалось 2, получено %v", values[1])
	}
	if values[2] != "ожидание" {
		t.Errorf("Неверный статус: ожидалось 'ожидание', получено '%v'", values[2])
	}
}

func TestInterviewReflectorMethods(t *testing.T) {
	now := time.Now()
	interview := &store.Interview{
		ApplicationID: 1,
		Date:          now,
		Result:        "запланировано",
	}

	if interview.GetNameDB() != "interviews" {
		t.Errorf("Неверное имя таблицы для собеседования: %s", interview.GetNameDB())
	}

	values := interview.GetValues()
	if values[0] != 1 {
		t.Errorf("Неверный ApplicationID: ожидалось 1, получено %v", values[0])
	}
	if values[1] != now {
		t.Errorf("Неверная дата: ожидалось %v, получено %v", now, values[1])
	}
	if values[2] != "запланировано" {
		t.Errorf("Неверный результат: ожидалось 'запланировано', получено '%v'", values[2])
	}

	fields := interview.GetFields()
	if len(fields) != 6 {
		t.Errorf("Ожидалось 6 полей для собеседования, получено %d", len(fields))
	}
}
