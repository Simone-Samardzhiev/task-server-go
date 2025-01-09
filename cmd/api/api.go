package api

import (
	"database/sql"
	"log"
	"net/http"
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
	log.Println("Starting API server...")
	mux := http.NewServeMux()
	return http.ListenAndServe(s.port, mux)
}
