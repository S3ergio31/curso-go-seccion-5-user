package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/S3ergio31/curso-go-seccion-5-user/internal/user"
	"github.com/S3ergio31/curso-go-seccion-5-user/pkg/bootstrap"
	"github.com/S3ergio31/curso-go-seccion-5-user/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db, err := bootstrap.DBConnection()
	logger := bootstrap.InitLogger()

	if err != nil {
		logger.Fatalln(err)
	}

	userRepository := user.NewRepository(logger, db)
	userService := user.NewService(userRepository, logger)
	userEndpoints := user.MakeEndpoints(userService)

	/*router.HandleFunc("/users", userEndpoints.Create).Methods("POST")
	router.HandleFunc("/users", userEndpoints.GetAll).Methods("GET")
	router.HandleFunc("/users/{id}", userEndpoints.Get).Methods("GET")
	router.HandleFunc("/users/{id}", userEndpoints.Update).Methods("PATCH")
	router.HandleFunc("/users/{id}", userEndpoints.Delete).Methods("DELETE")*/

	address := fmt.Sprintf("%s:%s", os.Getenv("APP_URL"), os.Getenv("APP_PORT"))
	server := &http.Server{
		Handler:      handler.NewUserHttpServer(userEndpoints),
		Addr:         address,
		WriteTimeout: 1 * time.Minute,
		ReadTimeout:  1 * time.Minute,
	}

	errCh := make(chan error)
	go func() {
		logger.Println("listen in ", address)
		errCh <- server.ListenAndServe()
	}()

	err = <-errCh

	if err != nil {
		logger.Fatal(err)
	}
}
