//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2020 SeMI Technologies B.V. All rights reserved.
//
//  CONTACT: hello@semi.technology
//

package esvector

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v5/esapi"
	"github.com/go-openapi/strfmt"
	"github.com/semi-technologies/weaviate/entities/filters"
	"github.com/semi-technologies/weaviate/entities/models"
	"github.com/semi-technologies/weaviate/entities/schema"
	"github.com/semi-technologies/weaviate/entities/schema/kind"
	"github.com/semi-technologies/weaviate/entities/search"
	"github.com/semi-technologies/weaviate/usecases/traverser"
	"github.com/sirupsen/logrus"
)

// ThingSearch searches for all things with optional filters without vector scoring
func (r *Repo) ThingSearch(ctx context.Context, limit int,
	filters *filters.LocalFilter, underscore traverser.UnderscoreProperties) (search.Results, error) {
	return r.search(ctx, allThingIndices, nil, limit, filters, traverser.GetParams{
		UnderscoreProperties: underscore,
	})
}

// ActionSearch searches for all things with optional filters without vector scoring
func (r *Repo) ActionSearch(ctx context.Context, limit int,
	filters *filters.LocalFilter, underscore traverser.UnderscoreProperties) (search.Results, error) {
	return r.search(ctx, allActionIndices, nil, limit, filters, traverser.GetParams{
		UnderscoreProperties: underscore,
	})
}

// ThingByID extracts the one result matching the ID. Returns nil on no results
// (without errors), but errors if it finds more than 1 results
func (r *Repo) ThingByID(ctx context.Context, id strfmt.UUID,
	params traverser.SelectProperties, underscore traverser.UnderscoreProperties) (*search.Result, error) {
	return r.searchByID(ctx, allThingIndices, id, params, underscore)
}

// ActionByID extracts the one result matching the ID. Returns nil on no results
// (without errors), but errors if it finds more than 1 results
func (r *Repo) ActionByID(ctx context.Context, id strfmt.UUID,
	params traverser.SelectProperties, underscore traverser.UnderscoreProperties) (*search.Result, error) {
	return r.searchByID(ctx, allActionIndices, id, params, underscore)
}

// Exists checks if an object with the id exists, if not, it forces a refresh
// and retries once.
func (r *Repo) Exists(ctx context.Context, id strfmt.UUID) (bool, error) {
	ok, err := r.exists(ctx, id)
	if err != nil {
		return false, fmt.Errorf("check exists: %v", err)
	}

	if ok {
		return true, nil
	}

	err = r.forceRefresh(ctx)
	if err != nil {
		return false, fmt.Errorf("check exists: force refresh: %v", err)
	}

	ok, err = r.exists(ctx, id)
	if err != nil {
		return false, fmt.Errorf("check exists: after forced refresh: %v", err)
	}

	return ok, nil
}

func (r *Repo) exists(ctx context.Context, id strfmt.UUID) (bool, error) {
	res, err := r.searchByID(ctx, allClassIndices, id, nil, traverser.UnderscoreProperties{})
	return res != nil, err
}

func (r *Repo) forceRefresh(ctx context.Context) error {
	req := esapi.IndicesRefreshRequest{
		Index: []string{allClassIndices},
	}

	res, err := req.Do(ctx, r.client)
	if err != nil {
		return fmt.Errorf("index refresh request: %v", err)
	}

	if err := errorResToErr(res, r.logger); err != nil {
		return fmt.Errorf("index refresh request: %v", err)
	}

	return nil
}

func (r *Repo) searchByID(ctx context.Context, index string, id strfmt.UUID,
	properties traverser.SelectProperties, underscore traverser.UnderscoreProperties) (*search.Result, error) {
	filters := &filters.LocalFilter{
		Root: &filters.Clause{
			On:       &filters.Path{Property: schema.PropertyName(keyID)},
			Value:    &filters.Value{Value: id},
			Operator: filters.OperatorEqual,
		},
	}
	res, err := r.search(ctx, index, nil, 2, filters, traverser.GetParams{
		Properties:           properties,
		UnderscoreProperties: underscore,
	})
	if err != nil {
		return nil, err
	}

	switch len(res) {
	case 0:
		return nil, nil
	case 1:
		return &res[0], nil
	default:
		return nil, fmt.Errorf("invalid number of results (%d) for id '%s'", len(res), id)
	}
}

type counterImpl struct {
	sync.Mutex
	count int
}

func (c *counterImpl) Inc() {
	c.Lock()
	defer c.Unlock()

	c.count++
}

func (c *counterImpl) Get() int {
	c.Lock()
	defer c.Unlock()

	return c.count
}

// ClassSearch searches for classes with optional filters without vector scoring
func (r *Repo) ClassSearch(ctx context.Context, params traverser.GetParams) ([]search.Result, error) {
	ctx, cancel := limitUnlimitedContext(ctx)
	defer cancel()

	start := time.Now()
	r.requestCounter = &counterImpl{}
	index := classIndexFromClassName(params.Kind, params.ClassName)
	res, err := r.search(ctx, index, nil, params.Pagination.Limit, params.Filters, params)
	count := r.requestCounter.(*counterImpl).Get()
	r.logger.WithFields(logrus.Fields{
		"action":        "esvector_class_search",
		"request_count": count,
		"took":          time.Since(start),
	}).Debug("completed class search")

	return res, err
}

// VectorClassSearch limits the vector search to a specific class (and kind)
func (r *Repo) VectorClassSearch(ctx context.Context, params traverser.GetParams) ([]search.Result, error) {
	ctx, cancel := limitUnlimitedContext(ctx)
	defer cancel()

	start := time.Now()
	r.requestCounter = &counterImpl{}
	index := classIndexFromClassName(params.Kind, params.ClassName)
	res, err := r.search(ctx, index, params.SearchVector, params.Pagination.Limit, params.Filters, params)
	count := r.requestCounter.(*counterImpl).Get()
	r.logger.WithFields(logrus.Fields{
		"action":        "esvector_vector_class_search",
		"request_count": count,
		"took":          time.Since(start),
	}).Debug("completed class search")

	return res, err
}

// VectorSearch retrives the closest concepts by vector distance
func (r *Repo) VectorSearch(ctx context.Context, vector []float32,
	limit int, filters *filters.LocalFilter) ([]search.Result, error) {
	return r.search(ctx, "*", vector, limit, filters, traverser.GetParams{})
}

func (r *Repo) search(ctx context.Context, index string,
	vector []float32, limit int,
	filters *filters.LocalFilter, params traverser.GetParams) ([]search.Result, error) {
	r.logger.
		WithField("action", "esvector_search").
		WithField("index", index).
		WithField("vector", vector).
		WithField("filters", filters).
		WithField("params", params).
		Debug("starting search to esvector")

	r.requestCounter.Inc()

	var buf bytes.Buffer

	query, err := r.queryFromFilter(ctx, filters)
	if err != nil {
		if _, ok := err.(SubQueryNoResultsErr); ok {
			// a sub-query error'd with no results, that's not an error case to us,
			// this simply means, we return no results to the user
			return nil, nil
		}
		return nil, err
	}

	body := r.buildSearchBody(query, vector, limit)

	err = json.NewEncoder(&buf).Encode(body)
	if err != nil {
		return nil, fmt.Errorf("vector search: encode json: %v", err)
	}

	res, err := r.client.Search(
		r.client.Search.WithContext(ctx),
		r.client.Search.WithIndex(index),
		r.client.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("vector search: %v", err)
	}

	return r.searchResponse(ctx, res, params.Properties, params.UnderscoreProperties)
}

func (r *Repo) buildSearchBody(filterQuery map[string]interface{}, vector []float32, limit int) map[string]interface{} {
	var query map[string]interface{}

	if vector == nil {
		query = filterQuery
	} else {
		query = map[string]interface{}{
			"function_score": map[string]interface{}{
				"query":      filterQuery,
				"boost_mode": "replace",
				"functions": []interface{}{
					map[string]interface{}{
						"script_score": map[string]interface{}{
							"script": map[string]interface{}{
								"inline": "binary_vector_score",
								"lang":   "knn",
								"params": map[string]interface{}{
									"cosine": true,
									"field":  keyVector,
									"vector": vector,
								},
							},
						},
					},
				},
			},
		}
	}

	return map[string]interface{}{
		"query": query,
		"size":  limit,
	}
}

type searchResponse struct {
	Hits struct {
		Hits []hit `json:"hits"`
	} `json:"hits"`
	Aggregations aggregations `json:"aggregations"`
}

type aggregations map[string]interface{}

// type singleAggregation struct {
// 	Buckets []map[string]interface{} `json:"buckets"`
// }

type hit struct {
	ID     string                 `json:"_id"`
	Source map[string]interface{} `json:"_source"`
	Score  float32                `json:"_score"`
	Index  string                 `json:"_index"`
}

func (r *Repo) searchResponse(ctx context.Context, res *esapi.Response,
	properties traverser.SelectProperties, underscoreProps traverser.UnderscoreProperties) ([]search.Result, error) {
	if err := errorResToErr(res, r.logger); err != nil {
		return nil, fmt.Errorf("vector search: %v", err)
	}

	var sr searchResponse
	defer res.Body.Close()
	err := json.NewDecoder(res.Body).Decode(&sr)
	if err != nil {
		return nil, fmt.Errorf("vector search: decode json: %v", err)
	}

	requestCacher := newCacher(r)
	err = requestCacher.build(ctx, sr, properties, underscoreProps.RefMeta)
	if err != nil {
		return nil, fmt.Errorf("build request cache: %v", err)
	}

	return sr.toResults(r, properties, underscoreProps, requestCacher)
}

func (sr searchResponse) toResults(r *Repo, properties traverser.SelectProperties,
	underscoreProps traverser.UnderscoreProperties, requestCacher *cacher) ([]search.Result, error) {
	hits := sr.Hits.Hits
	output := make([]search.Result, len(hits))

	for i, hit := range hits {
		var err error
		properties, err = requestCacher.replaceInitialPropertiesWithSpecific(hit, properties)
		if err != nil {
			return nil, err
		}

		k, err := kind.Parse(hit.Source[keyKind.String()].(string))
		if err != nil {
			return nil, fmt.Errorf("vector search: result %d: parse Kind: %v", i, err)
		}

		vector, err := base64ToVector(hit.Source[keyVector.String()].(string))
		if err != nil {
			return nil, fmt.Errorf("vector search: result %d: parse vector: %v", i, err)
		}

		schema, err := r.parseSchema(hit.Source, properties, underscoreProps, requestCacher)
		if err != nil {
			return nil, fmt.Errorf("vector search: result %d: parse schema: %v", i, err)
		}

		weights, err := parseVectorWeights(hit.Source[keyVectorWeights.String()])
		if err != nil {
			return nil, fmt.Errorf("vector search: result %d: parse weights: %v", i, err)
		}

		created := parseFloat64(hit.Source, keyCreated.String())
		updated := parseFloat64(hit.Source, keyUpdated.String())

		output[i] = search.Result{
			ClassName:     hit.Source[keyClassName.String()].(string),
			ID:            strfmt.UUID(hit.ID),
			Kind:          k,
			Score:         hit.Score,
			Vector:        vector,
			Schema:        schema,
			Created:       int64(created),
			Updated:       int64(updated),
			VectorWeights: weights,
		}
		if underscoreProps.Classification ||
			underscoreProps.Vector ||
			underscoreProps.Interpretation {
			underscores := &models.UnderscoreProperties{}
			storedUnderscores := r.extractUnderscoreProps(hit.Source)

			if underscoreProps.Vector {
				underscores.Vector = vector
			}

			if storedUnderscores != nil && underscoreProps.Classification {
				underscores.Classification = storedUnderscores.Classification
			}

			if storedUnderscores != nil && underscoreProps.Interpretation {
				underscores.Interpretation = storedUnderscores.Interpretation
			}

			output[i].UnderscoreProperties = underscores
		}
	}

	return output, nil
}

func parseFloat64(source map[string]interface{}, key string) float64 {
	untyped := source[key]
	if v, ok := untyped.(float64); ok {
		return v
	}

	return 0
}

func limitUnlimitedContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return context.WithCancel(ctx)
	}

	return context.WithTimeout(ctx, 15*time.Second)
}

func parseVectorWeights(in interface{}) (map[string]string, error) {
	if in == nil {
		return nil, nil
	}

	asMap, ok := in.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected map, got %T", asMap)
	}

	out := make(map[string]string, len(asMap))
	for key, value := range asMap {
		asString, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("key '%s': expected string, got %T", key, value)
		}

		out[key] = asString
	}

	return out, nil
}
