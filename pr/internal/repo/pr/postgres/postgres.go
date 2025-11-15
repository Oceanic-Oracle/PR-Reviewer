package pr_postgres

import (
	"context"
	"errors"
	"log/slog"
	apperrors "pr/internal/app/errors"
	"pr/internal/repo/pr"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PRPostgres struct {
	conn *pgxpool.Pool
	log  *slog.Logger
}

func (prp *PRPostgres) CreatePR(ctx context.Context, prm pr.PRModel) ([]string, error) {
	tx, err := prp.conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var authorTeam string
	err = tx.QueryRow(ctx, "SELECT team_name FROM users WHERE id = $1", prm.AuthorId).Scan(&authorTeam)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	rows, err := tx.Query(ctx, `
		SELECT id
		FROM users
		WHERE team_name = $1
		  AND id != $2
		  AND is_active = true
		ORDER BY RANDOM()
		LIMIT 2`,
		authorTeam, prm.AuthorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviewers []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, id)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO pull_requests (id, name, author_id, status)
		VALUES ($1, $2, $3, $4)`,
		prm.Id, prm.Name, prm.AuthorId, prm.Status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, apperrors.ErrPRExists
			}
		}
		return nil, err
	}

	for _, reviewerID := range reviewers {
		_, err := tx.Exec(ctx, `
			INSERT INTO users_pull_requests (pull_requests_id, users_id)
			VALUES ($1, $2)`,
			prm.Id, reviewerID)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return reviewers, nil
}

func (prp *PRPostgres) MergePR(ctx context.Context, id string) (*pr.PRWithDateModel, []string, error) {
	row := prp.conn.QueryRow(ctx, `
		UPDATE pull_requests 
		SET status = 'MERGED' 
		WHERE id = $1 
		RETURNING id, name, author_id, status, created_at, merged_at`,
		id)

	var prModel pr.PRWithDateModel
	err := row.Scan(
		&prModel.Id,
		&prModel.Name,
		&prModel.AuthorId,
		&prModel.Status,
		&prModel.CreatedAt,
		&prModel.MergedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, apperrors.ErrNotFound
		}
		return nil, nil, err
	}

	reviewers, err := prp.GetPrReviewers(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	return &prModel, reviewers, nil
}

func (prp *PRPostgres) GetPrReviewers(ctx context.Context, id string) ([]string, error) {
	var usersId []string

	sql := `
	SELECT u.id 
	FROM users AS u
		JOIN users_pull_requests AS upr
			ON u.id = upr.users_id
		JOIN pull_requests AS pr
			ON upr.pull_requests_id = pr.id
	WHERE pr.id = $1`

	rows, err := prp.conn.Query(ctx, sql, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var body string
		if err = rows.Scan(&body); err != nil {
			return nil, err
		}
		usersId = append(usersId, body)
	}

	return usersId, nil
}

func (prp *PRPostgres) SwapPRReviewer(ctx context.Context, prId, userId string) (pr.PRModel, []string, string, error) {
	tx, err := prp.conn.Begin(ctx)
	if err != nil {
		return pr.PRModel{}, nil, "", err
	}
	defer tx.Rollback(ctx)

	var userExists int
	err = tx.QueryRow(ctx, `
		SELECT 1
		FROM users
		WHERE id = $1`, userId).Scan(&userExists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pr.PRModel{}, nil, "", apperrors.ErrNotFound
		}
		return pr.PRModel{}, nil, "", err
	}

	var bodyPr pr.PRModel
	err = tx.QueryRow(ctx, `
		SELECT id, name, author_id, status
		FROM pull_requests
		WHERE id = $1`, prId).Scan(
		&bodyPr.Id, &bodyPr.Name, &bodyPr.AuthorId, &bodyPr.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pr.PRModel{}, nil, "", apperrors.ErrNotFound
		}
		return pr.PRModel{}, nil, "", err
	}

	if bodyPr.Status == "MERGED" {
		return pr.PRModel{}, nil, "", apperrors.ErrMerged
	}

	var assigned int
	err = tx.QueryRow(ctx, `
		SELECT 1
		FROM users_pull_requests
		WHERE pull_requests_id = $1 AND users_id = $2`,
		prId, userId).Scan(&assigned)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pr.PRModel{}, nil, "", apperrors.ErrNotAssigned
		}
		return pr.PRModel{}, nil, "", err
	}

	var teamName string
	err = tx.QueryRow(ctx, "SELECT team_name FROM users WHERE id = $1", bodyPr.AuthorId).Scan(&teamName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pr.PRModel{}, nil, "", apperrors.ErrNotFound
		}
		return pr.PRModel{}, nil, "", err
	}

	var newReviewerID string
	err = tx.QueryRow(ctx, `
		SELECT u.id
		FROM users u
		WHERE u.team_name = $1
		  AND u.is_active = true
		  AND u.id != $2
		  AND u.id != $3
		  AND u.id NOT IN (
		      SELECT users_id FROM users_pull_requests WHERE pull_requests_id = $4
			  	AND users_id != $2
		  )
		ORDER BY RANDOM()
		LIMIT 1`,
		teamName, userId, bodyPr.AuthorId, prId).Scan(&newReviewerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pr.PRModel{}, nil, "", apperrors.ErrNoCandidate
		}
		return pr.PRModel{}, nil, "", err
	}

	_, err = tx.Exec(ctx, `
		UPDATE users_pull_requests
		SET users_id = $1
		WHERE pull_requests_id = $2 AND users_id = $3`,
		newReviewerID, prId, userId)
	if err != nil {
		return pr.PRModel{}, nil, "", err
	}

	rows, err := tx.Query(ctx, `
		SELECT users_id FROM users_pull_requests WHERE pull_requests_id = $1`, prId)
	if err != nil {
		return pr.PRModel{}, nil, "", err
	}
	defer rows.Close()

	var newReviewers []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return pr.PRModel{}, nil, "", err
		}
		newReviewers = append(newReviewers, id)
	}

	if err := tx.Commit(ctx); err != nil {
		return pr.PRModel{}, nil, "", err
	}

	return bodyPr, newReviewers, newReviewerID, nil
}

func NewPrPostgres(conn *pgxpool.Pool, log *slog.Logger) *PRPostgres {
	return &PRPostgres{
		conn: conn,
		log:  log,
	}
}
