package api

import (
	"context"
	"github.com/khedhrije/podcaster-search-api/internal/domain/port"
	"github.com/khedhrije/podcaster-search-api/pkg"
)

type Search interface {
	ProgramByID(ctx context.Context, id string) ([]*pkg.ProgramResponse, error)
	Programs(ctx context.Context) ([]*pkg.ProgramResponse, error)
}

type searchApi struct {
	programSearchAdapter port.Program
}

func NewSearchApi(programSearchAdapter port.Program) Search {
	return &searchApi{
		programSearchAdapter: programSearchAdapter,
	}
}

func (api searchApi) ProgramByID(ctx context.Context, id string) ([]*pkg.ProgramResponse, error) {

	programs, err := api.programSearchAdapter.ProgramByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var response []*pkg.ProgramResponse
	for _, program := range programs {
		response = append(response, &pkg.ProgramResponse{
			ID:          program.ID,
			Name:        program.Name,
			Description: program.Description,
		})
	}

	return response, nil
}

func (api searchApi) Programs(ctx context.Context) ([]*pkg.ProgramResponse, error) {

	programs, err := api.programSearchAdapter.Programs(ctx)
	if err != nil {
		return nil, err
	}

	var response []*pkg.ProgramResponse
	for _, program := range programs {
		response = append(response, &pkg.ProgramResponse{
			ID:          program.ID,
			Name:        program.Name,
			Description: program.Description,
		})
	}

	return response, nil
}
