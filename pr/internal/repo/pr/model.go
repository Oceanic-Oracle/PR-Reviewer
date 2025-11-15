package pr

import "time"

type PRModel struct {
	Id       string
	Name     string
	AuthorId string
	Status   string
}

type PRWithDateModel struct {
	Id        string
	Name      string
	AuthorId  string
	Status    string
	CreatedAt *time.Time
	MergedAt  *time.Time
}
