package server

import (
	"CustomerCRUD/pkg/repository"

	"github.com/gorilla/mux"
)

type Server struct {
	Router     *mux.Router
	repository repository.CustomerRepository // TODO: Abstract service layer
}

func NewServer(repository repository.CustomerRepository) *Server {
	return &Server{
		repository: repository,
	}
}
