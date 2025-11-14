package user_postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	apperrors "pr/internal/app/errors"
	"pr/internal/repo/user"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	defer tx.Rollback(ctx)
	
	sqlTeam := `
	INSERT INTO teams (name)
	VALUES ($1)
	;
	`

	if _, err := tx.Exec(ctx, sqlTeam, team); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", nil, apperrors.ErrTeamExists
		} else {
			return "", nil, apperrors.ErrTeamExists
		}
	}

	if len(users) > 0 {	
		arrSqlUser := make([]string, 0, len(users))
		args := make([]interface{}, 0, len(users)*4)

		arrSqlUser = append(arrSqlUser,
		`INSERT INTO users (id, username, team_name, is_active) VALUES`)
		for ind, val := range users {
			var str string
			if ind == len(users) - 1 {
				str =  "($%d, $%d, $%d, $%d)"
			} else {
				str =  "($%d, $%d, $%d, $%d),"
			}
			arrSqlUser = append(arrSqlUser, fmt.Sprintf(str, ind*4+1, ind*4+2, ind*4+3, ind*4+4))
			args = append(args, val.Id, val.UserName, team, val.IsActive)
		}
		arrSqlUser = append(arrSqlUser, `
		ON CONFLICT (id)
		DO UPDATE SET
			username = EXCLUDED.username,
			team_name = EXCLUDED.team_name,
			is_active = EXCLUDED.is_active;`)
		
		if _, err := tx.Exec(ctx, strings.Join(arrSqlUser, "\n"), args...); err != nil {
			
			return "", nil, err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return "", nil, err
	}

	return team, users, nil
}

func (u *UserPostgres) GetUsersByTeam(ctx context.Context, team string) ([]user.UserModel, error) {
	sql := `SELECT id, username, team_name, is_active FROM teams LEFT JOIN users ON team_name = name WHERE team_name = $1`
	
	rows, err := u.conn.Query(ctx, sql, team)
	if err != nil {

	}
	defer rows.Close()

	var res []user.UserModel
	for rows.Next() {
		var body user.UserModel
		if err := rows.Scan(&body.Id, &body.UserName, &body.TeamName, &body.IsActive); err != nil {
			return nil, err
		}
		res = append(res, body)
	}

	if len(res) == 0 {
		return nil, apperrors.ErrNotFound
	}

	return res, nil
}

func (u *UserPostgres) SetActiveUser(ctx context.Context, id string, isActive bool) (user.UserModel, error) {
	sql := `UPDATE users SET is_active = $1 WHERE id = $2 RETURNING id, username, team_name, is_active`

	var res user.UserModel
	if err := u.conn.QueryRow(ctx, sql, id, isActive).Scan(&res.Id, &res.UserName, &res.TeamName, &res.IsActive); err != nil {
		if err == pgx.ErrNoRows {
			return user.UserModel{}, apperrors.ErrNotFound
		}
		return user.UserModel{}, err
	}

	return res, nil
}



func NewPrPostgres(conn *pgxpool.Pool, log *slog.Logger) user.UserInterface {
	return &UserPostgres{
		conn: conn,
		log: log,
	}
}