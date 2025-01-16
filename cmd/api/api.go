package api

import (
	"database/sql"
	"net/http"
	"server/auth"
	"server/user"
)

// Server is a struct used to create a new server instance
// that can be run.
type Server struct {
	// The port of the server.
	port string
	// The database connection.
	db *sql.DB
}

// NewServer will create a server with a port and connection to database.
func NewServer(address string, db *sql.DB) *Server {
	return &Server{
		port: address,
		db:   db,
	}
}

// Start will make the server lister to the port.
func (s *Server) Start() error {
	userRepository := user.NewPostgresRepository(s.db)
	userService := user.NewDefaultService(userRepository)
	userHandler := user.NewDefaultHandler(userService)

	mux := http.NewServeMux()
	mux.Handle("POST /users/register", userHandler.Register())
	mux.Handle("POST /users/login", userHandler.Login())
	mux.Handle("GET /users/refresh", auth.JWTMiddleware(userHandler.Refresh(), auth.RefreshToken))

	return http.ListenAndServe(s.port, mux)
}
