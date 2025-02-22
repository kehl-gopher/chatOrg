package data

type Company struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	ApiKey   string `json:"api_key,omitempty"`
	Password string `json:"-"`
}

type AIAnswer struct {
	Answer string `json:"answer"`
}

type About struct {
	ID        string `json:"id,omitempty"`
	About     string `json:"about"`
	Embedding string `json:"embedding"`
	CompanyID string `json:"company_id,omitempty"` // If needed
}

type Document struct {
	ID      string `json:"id,omitempty"`
	Content string `json:"content"`
	// DocumentPath string `json:"document_path"`
	Embedding string `json:"embedding"`
	CompanyID string `json:"company_id,omitempty"`
}
type Settings struct {
	Label    string `json:"label"`
	Type     string `json:"type"`
	Default  string `json:"default"`
	Required bool   `json:"required"`
}
type JsonData struct {
	ChannelID string     `json:"channel_id"`
	Settings  []Settings `json:"settings"`
	Message   string     `json:"message"`
}
