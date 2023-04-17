package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func APIServer(url string, port string) *Server {

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/allUsers", getUsers) // To request all users
	// r.HandleFunc("/user/{name}", user).Methods("GET")             // To request a specific user
	// r.HandleFunc("createUser", CreateUser).Methods("POST") // To create a new user
	// r.HandleFunc("deleteUser", DeleteUser).Methods("DELETE")      // To delete a user
	// r.HandleFunc("changePassword", ChangePassword).Methods("PUT") // To change a user's password
	// s := NewServer()
	// s.HandleFunc("/api/users", s.CreateUser()).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":10000", r))
	return nil
}

func corsHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		h(w, r)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	GetAllUsers(db)
	fmt.Println("Endpoint hit: returnAllGroceries")
}

// func getUsers(w http.ResponseWriter, r *http.Request) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		switch r.Method {
// 		case http.MethodGet:
// 			// gorm for db read / write
// 			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 			if err != nil {
// 				panic(err)
// 			}
// 			GetAllUsers(db)
// 		default:
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		}
// 	}
// }

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
