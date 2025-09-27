package handler

import (
	"context"
	"encoding/json"
	"net/http"

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
		endpoint.Endpoint(endpoints.Create), decodeCreateUser, encodeResponse, opts...,
	)).Methods("POST")

	return router
}

func decodeCreateUser(_ context.Context, r *http.Request) (any, error) {
	var request user.CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}

	return request, nil
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
