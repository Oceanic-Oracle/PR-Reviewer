package pr

import "time"

type PRModel struct {
	ID       string
	Name     string
	AuthorID string
	Status   string
}

type PRWithDateModel struct {
	ID        string
	Name      string
	AuthorID  string
	Status    string
	CreatedAt *time.Time
	MergedAt  *time.Time
}
