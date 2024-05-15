package elasticsearchv7

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/khedhrije/podcaster-search-api/internal/configuration"
	"net"
	"net/http"
	"time"
)

type SearchClient interface {
	RunQuery(ctx context.Context, indexName string, query map[string]interface{}) (string, error)
	RunQueryWithScroll(ctx context.Context, query map[string]interface{}) (string, error)
}

type client struct {
	es *elasticsearch.Client
}

// NewElasticSearchClient initializes a new Elasticsearch adapter with the given configuration.
func NewElasticSearchClient(config *configuration.AppConfig) (SearchClient, error) {
	esClient, err := createClient(config)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch adapter init failed: %w", err)
	}
	return &client{es: esClient}, nil
}

func (client client) RunQuery(ctx context.Context, indexName string, query map[string]interface{}) (string, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return "", err
	}

	res, err := client.es.Search(
		client.es.Search.WithContext(context.Background()),
		client.es.Search.WithIndex(indexName),
		client.es.Search.WithBody(&buf),
		client.es.Search.WithTrackTotalHits(true),
		client.es.Search.WithPretty(),
	)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return "", fmt.Errorf("error parsing the response body: %s", err)
		}
		// Print the response status and error information.
		return "", fmt.Errorf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return "", err
	}

	response, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(response), nil
}

func (client client) RunQueryWithScroll(ctx context.Context, query map[string]interface{}) (string, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return "", err
	}

	// Initial search request
	res, err := client.es.Search(
		client.es.Search.WithContext(context.Background()),
		client.es.Search.WithIndex("your-index-name"), // replace with your index name
		client.es.Search.WithBody(&buf),
		client.es.Search.WithScroll(time.Minute),
		client.es.Search.WithTrackTotalHits(true),
		client.es.Search.WithPretty(),
	)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return "", fmt.Errorf("error parsing the response body: %s", err)
		}
		return "", fmt.Errorf("[%s] %s: %s",
			res.Status(),
			e["error"].(map[string]interface{})["type"],
			e["error"].(map[string]interface{})["reason"],
		)
	}

	// Initial search response
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return "", err
	}

	scrollID := r["_scroll_id"].(string)

	// Subsequent scroll requests
	for {
		res, err = client.es.Scroll(
			client.es.Scroll.WithContext(context.Background()),
			client.es.Scroll.WithScrollID(scrollID),
			client.es.Scroll.WithScroll(time.Minute),
		)
		if err != nil {
			return "", err
		}
		defer res.Body.Close()

		if res.IsError() {
			var e map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
				return "", fmt.Errorf("error parsing the response body: %s", err)
			}
			return "", fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}

		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			return "", err
		}

		if len(r["hits"].(map[string]interface{})["hits"].([]interface{})) == 0 {
			break
		}

		scrollID = r["_scroll_id"].(string)
	}

	response, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	return string(response), nil
}

// createClient creates an Elasticsearch adapter using the given configuration.
func createClient(config *configuration.AppConfig) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{config.Elasticsearch.URL},
		Username:  config.Elasticsearch.User,
		Password:  config.Elasticsearch.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: 20 * time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS11,
				InsecureSkipVerify: true,
			},
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("elastic adapter creation failure: %w", err)
	}
	return client, nil
}
