package app

import (
	"context"

	api "github.com/MisLink/go-web-template/api"
)

func (*Server) Get(context.Context, *api.GetRequest) (*api.GetResponse, error) {
	return &api.GetResponse{}, nil
}
