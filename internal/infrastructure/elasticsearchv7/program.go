package elasticsearchv7

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/khedhrije/podcaster-search-api/internal/domain/model"
	"github.com/khedhrije/podcaster-search-api/internal/domain/port"
	"github.com/tidwall/gjson"
	"strings"
)

type programAdapter struct {
	client SearchClient
}

func NewProgramAdapter(client SearchClient) port.Program {
	return &programAdapter{
		client: client,
	}
}

func (adapter programAdapter) ProgramByID(ctx context.Context, id string) ([]*model.Program, error) {
	query := adapter.buildQueryForProgramByID(id)
	response, err := adapter.client.RunQuery(ctx, "latest", query)
	if err != nil {
		return nil, err
	}

	programs, err := extractPrograms(response)
	if err != nil {
		return nil, err
	}

	return programs, nil
}
func (adapter programAdapter) Programs(ctx context.Context) ([]*model.Program, error) {
	query := adapter.buildQueryForPrograms()
	response, err := adapter.client.RunQuery(ctx, "latest", query)
	if err != nil {
		return nil, err
	}

	programs, err := extractPrograms(response)
	if err != nil {
		return nil, err
	}

	return programs, nil
}

func (adapter programAdapter) buildQueryForPrograms() map[string]interface{} {

	must := make([]interface{}, 0)

	must = append(must, map[string]interface{}{
		"match_all": map[string]interface{}{},
	})

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
	}
	return query
}

func (adapter programAdapter) buildQueryForProgramByID(id string) map[string]interface{} {

	must := make([]interface{}, 0)

	must = append(must, map[string]interface{}{
		"term": map[string]interface{}{
			"ID": id,
		},
	})

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
	}
	return query
}

func extractPrograms(searchResult string) ([]*model.Program, error) {
	source := gjson.Get(searchResult, "hits.hits.#._source")
	locs := make([]*model.Program, 0)
	if err := json.NewDecoder(strings.NewReader(source.Raw)).Decode(&locs); err != nil {
		return nil, fmt.Errorf("json.NewDecoder.Decode : %w", err)
	}
	return locs, nil
}
