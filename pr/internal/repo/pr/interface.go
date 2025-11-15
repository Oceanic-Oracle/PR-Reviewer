package pr

import (
	"context"
)

type PRInterface interface {
	CreatePR(ctx context.Context, prm PRModel) ([]string, error)
	MergePR(ctx context.Context, id string) (*PRWithDateModel, []string, error)
	GetPrReviewers(ctx context.Context, id string) ([]string, error)
	SwapPRReviewer(ctx context.Context, prId, userId string) (PRModel, []string, string, error)
}
