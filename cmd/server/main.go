package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"PEAK/internal/db"
	"PEAK/internal/handlers"
	"PEAK/internal/middleware"

	"github.com/gorilla/mux"
)

func main() {
	if err := db.Connect(); err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/register", handlers.Register).Methods("POST")
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")

	// protected
	s := r.PathPrefix("/api").Subrouter()
	s.Use(middleware.JWTAuth)
	s.HandleFunc("/projects", handlers.CreateProject).Methods("POST")
	s.HandleFunc("/projects", handlers.ListProjects).Methods("GET")
	s.HandleFunc("/projects/{id}", handlers.GetProject).Methods("GET")

	s.HandleFunc("/defects", handlers.CreateDefect).Methods("POST")
	s.HandleFunc("/defects", handlers.ListDefects).Methods("GET")
	s.HandleFunc("/defects/{id}/status", handlers.UpdateDefectStatus).Methods("PATCH")

	srv := &http.Server{
		Handler: r,
		Addr: ":" + func() string {
			if p := os.Getenv("PORT"); p != "" {
				return p
			}
			return "5432"
		}(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
