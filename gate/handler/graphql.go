package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ghosv/open/meta"

	"github.com/ddosakura/sola/v2"
	"github.com/ghosv/open/gate/client"
	"github.com/ghosv/open/gate/graphql/loader"
	"github.com/ghosv/open/gate/graphql/resolver"
	"github.com/ghosv/open/gate/graphql/schema"

	"github.com/graph-gophers/graphql-go"
)

// GraphQL handler
type GraphQL struct {
	Schema            *graphql.Schema
	LoadersInitialize func(*client.MicroClient) loader.Collection
}

// Handler of GraphQL
func (h *GraphQL) Handler(c sola.Context) error {
	r := c.Request()

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		return err
	}

	service := c.Get(meta.CtxService).(*client.MicroClient)
	ctx := h.LoadersInitialize(service).Attach(r.Context())
	ctx = context.WithValue(ctx, meta.KeyService, service)
	ctx = context.WithValue(ctx, meta.KeyTokenPayload, c.Get(meta.CtxTokenPayload))

	response := h.Schema.Exec(ctx, params.Query, params.OperationName, params.Variables)

	return c.JSON(http.StatusOK, response)
}

// GraphQL
var (
	rootResolver = &resolver.Resolver{}

	// DefaultGraphQL POST /graphql <Need Login>
	DefaultGraphQL = (&GraphQL{
		Schema: graphql.MustParseSchema(schema.String(),
			rootResolver),
		LoadersInitialize: loader.Initialize,
	}).Handler
)
