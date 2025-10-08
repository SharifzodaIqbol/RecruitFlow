package helper

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
