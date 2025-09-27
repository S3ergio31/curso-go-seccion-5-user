package user

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/S3ergio31/curso-go-seccion-5-meta/meta"
	"github.com/S3ergio31/curso-go-seccion-5-response/response"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
)

type Controller func(w http.ResponseWriter, r *http.Request)

type Endpoints struct {
	Create endpoint.Endpoint
	Get    Controller
	GetAll Controller
	Update Controller
	Delete Controller
}

type CreateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type UpdateRequest struct {
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

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		user, err := s.Get(id)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: user})
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		filters := Filters{
			FirstName: query.Get("first_name"),
			LastName:  query.Get("last_name"),
		}

		count, err := s.Count(filters)

		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(Response{Status: 500, Err: err.Error()})
			return
		}

		limit, _ := strconv.Atoi(query.Get("limit"))
		page, _ := strconv.Atoi(query.Get("page"))
		meta, err := meta.New(page, limit, count, os.Getenv("PAGINATOR_LIMIT_DEFAULT"))

		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(Response{Status: 500, Err: err.Error()})
			return
		}

		users, err := s.GetAll(filters, meta.Offset(), meta.Limit())

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(Response{
			Status: 200,
			Data:   users,
			Meta:   meta,
		})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var updateRequest UpdateRequest

		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "invalid request format"})
			return
		}

		if updateRequest.FirstName != nil && *updateRequest.FirstName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "first name is required"})
			return
		}

		if updateRequest.LastName != nil && *updateRequest.LastName == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(Response{Status: 400, Err: "last name is required"})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		err := s.Update(
			id,
			updateRequest.FirstName,
			updateRequest.LastName,
			updateRequest.Email,
			updateRequest.Phone,
		)

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(Response{Status: 404, Err: "user does not exist"})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: "success"})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		err := s.Delete(id)

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(Response{Status: 404, Err: "user does not exists"})
			return
		}

		json.NewEncoder(w).Encode(Response{Status: 200, Data: "success"})
	}
}
