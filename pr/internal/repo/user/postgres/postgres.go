package user_postgres

import (
	"context"
	"fmt"
	"log/slog"
	"pr/internal/repo/user"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPostgres struct {
	conn *pgxpool.Pool
	log *slog.Logger
}

func (u *UserPostgres) CreateOrUpdateTeamWithUsers(ctx context.Context, team string,
		users []user.UserModel) (string, []user.UserModel, error) {
	tx, err := u.conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return "", nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	
	sqlTeam := `
	INSERT INTO teams (name)
	VALUES ($1)
	;
	`

	if _, err := tx.Exec(ctx, sqlTeam, team); err != nil {
		return "", nil, fmt.Errorf("team_name already exists")
	}

	arrSqlUser := make([]string, 0, len(users))
	args := make([]interface{}, 0, len(users)*4)

	arrSqlUser = append(arrSqlUser, `INSERT INTO users (id, username, team_name, is_active)`)
	for ind, val := range users {
		arrSqlUser = append(arrSqlUser, fmt.Sprintf(`VALUES ($%d, $%d, $%d, $%d)`, ind*4+1, ind*4+2, ind*4+3, ind*4+4))
		args = append(args, val.Id, val.UserName, team, val.IsActive)
	}
	arrSqlUser = append(arrSqlUser, `
	ON CONFLICT (id)
	DO UPDATE SET
    	username = EXCLUDED.username,
    	team_name = EXCLUDED.team_name,
    	is_active = EXCLUDED.is_active;`)
	
	if _, err := tx.Exec(ctx, strings.Join(arrSqlUser, "\n"), args); err != nil {
		return "", nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return "", nil, err
	}

	return team, users, nil
}

func NewPrPostgres(conn *pgxpool.Pool, log *slog.Logger) user.UserInterface {
	return &UserPostgres{
		conn: conn,
		log: log,
	}
}