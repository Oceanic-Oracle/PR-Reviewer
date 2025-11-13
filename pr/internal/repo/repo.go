package repo

import (
	"pr/internal/repo/pr"
	"pr/internal/repo/user"
)

type Repo struct {
	User user.UserInterface
	PR   pr.PRInterface
}

func NewRepo(userDb user.UserInterface, prDb pr.PRInterface) *Repo {
	return &Repo{
		User: userDb,
		PR:   prDb,
	}
}
