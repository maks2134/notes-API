package main

import (
	"database/sql"
	"log"
	"net/http"
	"notes-api/internal/api/handler"
	"notes-api/internal/config"
	"notes-api/internal/repository"
	"notes-api/internal/service"

	_ "notes-api/docs"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Notes API
// @version 1.0
// @description This is a sample Notes API with JWT authentication
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DB.GetPostgresDSN())
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer db.Close()

	noteRepo := repository.NewPostgresNoteRepository(db)
	userRepo := repository.NewUserRepository(db)

	noteService := service.NewNoteService(noteRepo)
	authService := service.NewAuthService(userRepo)

	noteHandler := handler.NewNoteHandler(noteService)
	authHandler := handler.NewAuthHandler(authService)

	r := mux.NewRouter()

	authHandler.RegisterRoutes(r)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(handler.AuthMiddleware)
	api.HandleFunc("/notes", noteHandler.GetNotes).Methods("GET")
	api.HandleFunc("/notes", noteHandler.CreateNote).Methods("POST")
	api.HandleFunc("/notes/{id}", noteHandler.GetNote).Methods("GET")
	api.HandleFunc("/notes/{id}", noteHandler.UpdateNote).Methods("PUT")
	api.HandleFunc("/notes/{id}", noteHandler.DeleteNote).Methods("DELETE")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Printf("Server started on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
