package entity

type Appeal struct {
	ID        uint
	UserID    uint
	Tag       string
	Name      string
	Status    string
	CreatedAt string
	UpdatedAt string
}

type AppealMessage struct {
	ID uint
	// Sender    string // "user" or "support"
	IsResponse bool
	Message    string
	CreatedAt  string
}
