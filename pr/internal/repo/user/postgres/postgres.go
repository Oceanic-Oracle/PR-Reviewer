package user_postgres

import (
	"context"
	"errors"
	"log/slog"
	apperrors "pr/internal/app/errors"
	"pr/internal/repo/pr"
	"pr/internal/repo/user"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPostgres struct {
	conn *pgxpool.Pool
	log *slog.Logger
}

func (u *UserPostgres) SetUserStatus(ctx context.Context, id string, isActive bool) (user.UserModel, error) {
	sql := `UPDATE users SET is_active = $1 WHERE id = $2 RETURNING id, username, team_name, is_active`

	var res user.UserModel
	if err := u.conn.QueryRow(ctx, sql, isActive, id).Scan(&res.Id, &res.UserName, &res.TeamName, &res.IsActive); err != nil {
		if err == pgx.ErrNoRows {
			return user.UserModel{}, apperrors.ErrNotFound
		}
		return user.UserModel{}, err
	}

	return res, nil
}

func (u *UserPostgres) GetPrUser(ctx context.Context, userId string) ([]pr.PRModel, error) {
	sql := `
	SELECT pr.id, pr.name, pr.author_id, pr.status
	FROM users_pull_requests AS upr
		JOIN pull_requests AS pr
			ON upr.pull_requests_id = pr.id
	WHERE upr.users_id = $1`
	rows, err := u.conn.Query(ctx, sql, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	defer rows.Close()

	var prs []pr.PRModel
	for rows.Next() {
		var body pr.PRModel
		if err := rows.Scan(&body.Id, &body.Name, &body.AuthorId, &body.Status); err != nil {
			return nil, err
		}
		prs = append(prs, body)
	}
	return prs, nil
}

func NewPrPostgres(conn *pgxpool.Pool, log *slog.Logger) user.UserInterface {
	return &UserPostgres{
		conn: conn,
		log: log,
	}
}