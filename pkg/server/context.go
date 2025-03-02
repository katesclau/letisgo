package server

import (
	"context"

	"mnesis.com/pkg/server/endpoints"
)

func getContext(api endpoints.Endpoints) (context.Context, context.CancelFunc) {
	ctx := context.WithValue(
		context.WithValue(
			context.WithValue(
				context.Background(),
				"api_name", api.Name,
			),
			"api_version", api.Version,
		),
		"api_description",
		api.Description,
	)
	return context.WithCancel(ctx)
}
