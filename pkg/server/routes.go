package server

import (
	"github.com/gorilla/mux"
)

func (s *Server) SetupRoutes() {
	s.Router = mux.NewRouter()

	s.Router.HandleFunc("/customers", s.GetAllCustomers).Methods("GET")
	s.Router.HandleFunc("/customers", s.CreateCustomer).Methods("POST")

	s.Router.HandleFunc("/customers/{id}", s.GetCustomerByID).Methods("GET")
	s.Router.HandleFunc("/customers/{id}", s.UpdateCustomer).Methods("PUT")
	s.Router.HandleFunc("/customers/{id}", s.DeleteCustomer).Methods("DELETE")

	s.Router.HandleFunc("/customers/email/{email}", s.GetCustomerByEmail).Methods("GET")
}
