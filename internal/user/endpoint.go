package user

import (
	"context"
	"net/http"
	"os"

	"github.com/S3ergio31/curso-go-seccion-5-meta/meta"
	"github.com/S3ergio31/curso-go-seccion-5-response/response"
	"github.com/go-kit/kit/endpoint"
)

type Controller func(w http.ResponseWriter, r *http.Request)

type Endpoints struct {
	Create endpoint.Endpoint
	Get    endpoint.Endpoint
	GetAll endpoint.Endpoint
	Update endpoint.Endpoint
	Delete endpoint.Endpoint
}

type CreateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type GetRequest struct {
	ID string
}

type GetAllRequest struct {
	FirstName string
	LastName  string
	Limit     int
	Page      int
}

type DeleteRequest struct {
	ID string
}

type UpdateRequest struct {
	ID        string
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
}

type Response struct {
	Status int        `json:"status"`
	Data   any        `json:"data,omitempty"`
	Err    string     `json:"error,omitempty"`
	Meta   *meta.Meta `json:"meta,omitempty"`
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		createRequest := request.(CreateRequest)

		if createRequest.FirstName == "" {
			return nil, response.BadRequest("first name is required")
		}

		if createRequest.LastName == "" {
			return nil, response.BadRequest("last name is required")
		}

		user, err := s.Create(
			createRequest.FirstName,
			createRequest.LastName,
			createRequest.Email,
			createRequest.Phone,
		)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", user, nil), nil
	}
}

func makeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetRequest)
		user, err := s.Get(req.ID)

		if err != nil {
			return nil, response.NotFound(err.Error())
		}

		return response.Ok("success", user, nil), nil
	}
}

func makeGetAllEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetAllRequest)
		filters := Filters{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		count, err := s.Count(filters)

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, os.Getenv("PAGINATOR_LIMIT_DEFAULT"))

		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		users, err := s.GetAll(filters, meta.Offset(), meta.Limit())

		if err != nil {
			return nil, response.BadRequest(err.Error())
		}

		return response.Ok("success", users, nil), nil
	}
}

func makeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		updateRequest := request.(UpdateRequest)

		if updateRequest.FirstName != nil && *updateRequest.FirstName == "" {
			return nil, response.BadRequest("first name is required")
		}

		if updateRequest.LastName != nil && *updateRequest.LastName == "" {
			return nil, response.BadRequest("last name is required")
		}

		err := s.Update(
			updateRequest.ID,
			updateRequest.FirstName,
			updateRequest.LastName,
			updateRequest.Email,
			updateRequest.Phone,
		)

		if err != nil {
			return nil, response.BadRequest("user does not exist")
		}

		return response.Ok("success", nil, nil), nil
	}
}

func makeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(DeleteRequest)
		err := s.Delete(req.ID)

		if err != nil {
			return nil, response.NotFound("user does not exists")
		}

		return response.Ok("success", nil, nil), nil
	}
}
