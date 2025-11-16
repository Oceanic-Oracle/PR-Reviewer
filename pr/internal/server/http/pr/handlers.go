package pr

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	apperrors "pr/internal/app/errors"
	"pr/internal/dto"
	"pr/internal/logctx"
	"pr/internal/repo"
	"pr/internal/repo/pr"
	errorhandler "pr/internal/server/http/error"
)

func CreatePR(repos *repo.Repo, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logctx.WithReqID(r.Context())
		log.InfoContext(ctx, "Received CreatePR request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		var body dto.CreatePRRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.WarnContext(ctx, "Invalid JSON in CreatePR request", slog.Any("error", err))
			errorhandler.WriteError(w, apperrors.ErrBadRequest, log)

			return
		}

		if body.ID == "" || body.Name == "" || body.AuthorID == "" {
			log.WarnContext(ctx, "Missing required fields in CreatePR request")
			errorhandler.WriteError(w, apperrors.ErrBadRequest, log)

			return
		}

		ctx = logctx.WithPR(ctx, body.ID)
		ctx = logctx.WithUserID(ctx, body.AuthorID)

		prModel := pr.PRModel{
			ID:       body.ID,
			Name:     body.Name,
			AuthorID: body.AuthorID,
			Status:   "OPEN",
		}

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		reviewers, err := repos.PR.CreatePR(ctx, prModel)
		if err != nil {
			log.ErrorContext(ctx, "Failed to create PR", slog.Any("error", err))
			errorhandler.WriteError(w, err, log)

			return
		}

		response := dto.CreatePRResponse{
			PR: dto.PRWithReviewers{
				PR: dto.PR{
					ID:       body.ID,
					Name:     body.Name,
					AuthorID: body.AuthorID,
					Status:   "OPEN",
				},
				AssignedReviewers: reviewers,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.ErrorContext(ctx, "Failed to encode CreatePR response", slog.Any("error", err))
		}

		log.InfoContext(ctx, "PR created successfully")
	}
}

func MergePR(repos *repo.Repo, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logctx.WithReqID(r.Context())
		log.InfoContext(ctx, "Received MergePR request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		var body dto.MergePRRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.WarnContext(ctx, "Invalid JSON in MergePR request", slog.Any("error", err))
			errorhandler.WriteError(w, apperrors.ErrBadRequest, log)

			return
		}

		if body.ID == "" {
			log.WarnContext(ctx, "Missing pull_request_id in MergePR request")
			errorhandler.WriteError(w, apperrors.ErrBadRequest, log)

			return
		}

		ctx = logctx.WithPR(ctx, body.ID)

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		prmodel, reviewers, err := repos.PR.MergePR(ctx, body.ID)
		if err != nil {
			log.ErrorContext(ctx, "Failed to merge PR", slog.Any("error", err))
			errorhandler.WriteError(w, err, log)

			return
		}

		response := dto.MergePRResponse{
			PR: dto.PRWithReviewersAndMerge{
				PR: dto.PR{
					ID:       prmodel.ID,
					Name:     prmodel.Name,
					AuthorID: prmodel.AuthorID,
					Status:   prmodel.Status,
				},
				AssignedReviewers: reviewers,
				MergedAt:          prmodel.MergedAt,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.ErrorContext(ctx, "Failed to encode MergePR response", slog.Any("error", err))
		}

		log.InfoContext(ctx, "PR merged successfully")
	}
}

func SwapPRReviewer(repos *repo.Repo, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logctx.WithReqID(r.Context())
		log.InfoContext(ctx, "Received SwapPRReviewer request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		var body dto.SwapPRReviewerRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.WarnContext(ctx, "Invalid JSON in SwapPRReviewer request", slog.Any("error", err))
			errorhandler.WriteError(w, apperrors.ErrBadRequest, log)

			return
		}

		if body.PRID == "" || body.OldUserID == "" {
			log.WarnContext(ctx, "Missing required fields in SwapPRReviewer request",
				slog.String("pr_id", body.PRID),
				slog.String("old_user_id", body.OldUserID))
			errorhandler.WriteError(w, apperrors.ErrBadRequest, log)

			return
		}

		ctx = logctx.WithPR(ctx, body.PRID)
		ctx = logctx.WithUserID(ctx, body.OldUserID)

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		prm, reviewers, reviewer, err := repos.PR.SwapPRReviewer(ctx, body.PRID, body.OldUserID)
		if err != nil {
			log.ErrorContext(ctx, "Failed to swap PR reviewer", slog.Any("error", err))
			errorhandler.WriteError(w, err, log)

			return
		}

		response := dto.SwapPRReviewerResponse{
			PR: dto.PRWithReviewers{
				PR: dto.PR{
					ID:       prm.ID,
					Name:     prm.Name,
					Status:   prm.Status,
					AuthorID: prm.AuthorID,
				},
				AssignedReviewers: reviewers,
			},
			ReplacedBy: reviewer,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.ErrorContext(ctx, "Failed to encode SwapPRReviewer response", slog.Any("error", err))
		}

		log.InfoContext(ctx, "PR reviewer swapped successfully",
			slog.String("new_reviewer_id", reviewer))
	}
}
