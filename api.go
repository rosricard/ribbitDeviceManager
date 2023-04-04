package main

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

type user struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:name`
	Password string    `json:password`
}

type Server struct {
	*mux.Router

	users []user
}

func NewServer() *Server {
	s := &Server{
		Router: mux.NewRouter(),
		users:  []user{},
	}
	return s
}

func (s *Server) CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var i user
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		i.ID = uuid.Must(uuid.NewV4())
		s.users = append(s.users, i)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(i); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// TODO: func delete user
// TODO: func get all users
// TODO: change password
// TODO: get Goliath API key
