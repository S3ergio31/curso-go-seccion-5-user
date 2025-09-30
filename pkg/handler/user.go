package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/S3ergio31/curso-go-seccion-5-response/response"
	"github.com/S3ergio31/curso-go-seccion-5-user/internal/user"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewUserHttpServer(endpoints user.Endpoints) http.Handler {
	router := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateUser,
		encodeResponse,
		opts...,
	)).Methods("POST")

	router.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllUser,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetUser,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteUser,
		encodeResponse,
		opts...,
	)).Methods("DELETE")

	router.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateUser,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	return router
}

func decodeCreateUser(_ context.Context, r *http.Request) (any, error) {
	var request user.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	return request, nil
}

func decodeUpdateUser(_ context.Context, r *http.Request) (any, error) {
	var request user.UpdateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	path := mux.Vars(r)

	request.ID = path["id"]

	return request, nil
}

func decodeGetUser(_ context.Context, r *http.Request) (any, error) {
	p := mux.Vars(r)
	req := user.GetRequest{ID: p["id"]}

	return req, nil
}

func decodeDeleteUser(_ context.Context, r *http.Request) (any, error) {
	p := mux.Vars(r)
	req := user.DeleteRequest{ID: p["id"]}

	return req, nil
}

func decodeGetAllUser(_ context.Context, r *http.Request) (any, error) {
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := user.GetAllRequest{
		FirstName: v.Get("first_name"),
		LastName:  v.Get("last_name"),
		Limit:     limit,
		Page:      page,
	}

	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, res any) error {
	r := res.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())

	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
