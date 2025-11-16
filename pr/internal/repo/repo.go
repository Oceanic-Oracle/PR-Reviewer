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

func NewRepo(teamDB team.UserInterface, userDB user.UserInterface, prDB pr.PRInterface) *Repo {
	return &Repo{
		Team: teamDB,
		User: userDB,
		PR:   prDB,
	}
}
