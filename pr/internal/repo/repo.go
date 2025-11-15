package repo

import (
	"pr/internal/repo/pr"
	"pr/internal/repo/team"
	"pr/internal/repo/user"
)

type Repo struct {
	Team team.UserInterface
	User user.UserInterface
	PR   pr.PRInterface
}

func NewRepo(teamDb team.UserInterface, userDb user.UserInterface, prDb pr.PRInterface) *Repo {
	return &Repo{
		Team: teamDb,
		User: userDb,
		PR:   prDb,
	}
}
