package pr

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
	CreatedAt string
	MergedAt  string
}
