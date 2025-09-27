package user

import (
	"log"

	"github.com/S3ergio31/curso-go-seccion-5-domain/domain"
)

type Service interface {
	Create(firstName, lastName, email, phone string) (*domain.User, error)
	GetAll(filters Filters, offset, limit int) ([]domain.User, error)
	Get(id string) (*domain.User, error)
	Delete(id string) error
	Update(id string, firstName, lastName, email, phone *string) error
	Count(filters Filters) (int, error)
}

type Filters struct {
	FirstName string
	LastName  string
}

type service struct {
	logger     *log.Logger
	repository Repository
}

func (s service) Create(firstName, lastName, email, phone string) (*domain.User, error) {
	user := &domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	if err := s.repository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s service) GetAll(filters Filters, offset, limit int) ([]domain.User, error) {
	users, err := s.repository.GetAll(filters, offset, limit)

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s service) Get(id string) (*domain.User, error) {
	user, err := s.repository.Get(id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s service) Delete(id string) error {
	return s.repository.Delete(id)
}

func (s service) Update(id string, firstName, lastName, email, phone *string) error {
	return s.repository.Update(id, firstName, lastName, email, phone)
}

func (s service) Count(filters Filters) (int, error) {
	return s.repository.Count(filters)
}

func NewService(repository Repository, logger *log.Logger) Service {
	return &service{logger: logger, repository: repository}
}
