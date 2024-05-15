package port

import (
	"context"
	"github.com/khedhrije/podcaster-search-api/internal/domain/model"
)

type Program interface {
	ProgramByID(ctx context.Context, id string) ([]*model.Program, error)
	Programs(ctx context.Context) ([]*model.Program, error)
}
