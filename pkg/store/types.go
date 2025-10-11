package store

import (
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
	New() Reflector
}

func (c *Company) GetNameDB() string        { return "companies" }
func (c *Company) GetParam() string         { return "name, created_at, updated_at" }
func (c *Company) GetPlaceholder() string   { return "$1, NOW(), NOW()" }
func (c *Company) GetValues() []interface{} { return []interface{}{c.Name} }
func (c *Company) GetFields() []interface{} {
	return []interface{}{
		&c.ID,
		&c.Name,
		&c.CreatedAt,
		&c.UpdatedAt,
	}
}
func (c *Company) New() Reflector {
	return &Company{}
}

func (c *Candidate) GetNameDB() string { return "candidate" }
func (c *Candidate) GetParam() string {
	return "name, tel_number, email, created_at, updated_at"
}
func (c *Candidate) GetPlaceholder() string { return "$1, $2, $3, NOW(), NOW()" }
func (c *Candidate) GetValues() []interface{} {
	return []interface{}{c.Name, c.TelNumber, c.Email}
}
func (c *Candidate) GetFields() []interface{} {
	return []interface{}{
		&c.ID,
		&c.Name,
		&c.TelNumber,
		&c.Email,
		&c.CreatedAt,
		&c.UpdatedAt,
	}
}
func (c *Candidate) New() Reflector {
	return &Candidate{}
}

func (p *Posting) GetNameDB() string { return "job_postings" }
func (p *Posting) GetParam() string {
	return "company_id, title, description, created_at, updated_at"
}
func (p *Posting) GetPlaceholder() string { return "$1, $2, $3, NOW(), NOW()" }
func (p *Posting) GetValues() []interface{} {
	return []interface{}{p.CompanyID, p.Title, p.Description}
}
func (p *Posting) GetFields() []interface{} {
	return []interface{}{
		&p.ID,
		&p.CompanyID,
		&p.Title,
		&p.Description,
		&p.CreatedAt,
		&p.UpdatedAt,
	}
}
func (p *Posting) New() Reflector {
	return &Posting{}
}
