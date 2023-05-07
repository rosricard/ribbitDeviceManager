package api

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/rosricard/userAccess/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var dsn = "root:password@tcp(127.0.0.1:3306)/ribbit?charset=utf8mb4&parseTime=True&loc=Local"

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
	s.routes()
	return s
}

func (s *Server) routes() {
	s.HandleFunc("/allUsers", s.GetUsers()).Methods("GET") // To request all users
	// r.HandleFunc("/user/{name}", user).Methods("GET")             // To request a specific user
	// r.HandleFunc("createUser", CreateUser).Methods("POST") // To create a new user
	// r.HandleFunc("deleteUser", DeleteUser).Methods("DELETE")      // To delete a user
	// r.HandleFunc("changePassword", ChangePassword).Methods("PUT") // To change a user's password
	// s := NewServer()
	// s.HandleFunc("/api/users", s.CreateUser()).Methods(http.MethodPost)
}

func (s *Server) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		db.GetAllUsers(dbConn)
		fmt.Println("Endpoint hit: returnAllUsers")
	}
}

// func CreateUser(w http.ResponseWriter, r *http.Request) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var i user
// 		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		i.ID = uuid.Must(uuid.NewV4())
// 		s.users = append(s.users, i)

// 		w.Header().Set("Content-Type", "application/json")
// 		if err := json.NewEncoder(w).Encode(i); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}
// }

// TODO: func delete user
// TODO: func get all users
// TODO: func get a specific user
// TODO: change password
// TODO: get Goliath API key
