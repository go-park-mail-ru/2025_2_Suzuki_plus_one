package entity

type Appeal struct {
	ID        uint
	Tag       string
	Name      string
	Status    string
	CreatedAt string
	UpdatedAt string
}

type AppealMessage struct {
	ID        uint
	Sender    string // "user" or "support"
	Message   string
	CreatedAt string
}
