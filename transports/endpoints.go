package transports

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"parserProject/models"
	"parserProject/service"
)

var (
	NoaAuth    = errors.New("required: x-api-key")
	NotAllowed = errors.New("method not allowed")
	NotFound   = errors.New("page not found")
)

type Endpoints struct {
	GetTitleFromUrls endpoint.Endpoint
}

func MakeEndpoints(parser service.ServiceInterface, apiKey string) Endpoints {
	return Endpoints{
		GetTitleFromUrls: endpoint.Chain(middlewareApiKey(apiKey))(makeGetTitleFromUrlsEndpoint(parser)),
	}
}

func middlewareApiKey(apiKey string) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if apiKey != ctx.Value("x-api-key").(string) {
				return nil, NoaAuth
			} else {
				return next(ctx, request)
			}
		}
	}
}

func makeGetTitleFromUrlsEndpoint(parser service.ServiceInterface) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		requestBody := request.(models.UrlsRequest)
		res, err := parser.GetTitlesWithUrl(requestBody.Urls)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}
