package main

import (
	"database/sql"
	"log"
	"net/http"
	"notes-api/internal/api/handler"
	"notes-api/internal/config"
	"notes-api/internal/repository"
	"notes-api/internal/service"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	_ "notes-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Notes API
// @version 1.1
// @description Это сервер для приложения заметок с поддержкой чек-листов.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description "Type 'Bearer' followed by a space and JWT token."

func main() {
	log.Println("Загрузка конфигурации...")
	cfg := config.LoadConfig()

	log.Println("Подключение к базе данных...")
	db, err := sql.Open("postgres", cfg.DB.GetPostgresDSN())
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Не удалось проверить подключение к базе данных: %v", err)
	}
	log.Println("База данных успешно подключена.")

	userRepo := repository.NewUserRepository(db)
	noteRepo := repository.NewPostgresNoteRepository(db)
	checklistItemRepo := repository.NewPostgresChecklistItemRepository(db)

	authService := service.NewAuthService(userRepo)
	noteService := service.NewNoteService(noteRepo)
	checklistItemService := service.NewChecklistItemService(checklistItemRepo, noteRepo)

	authHandler := handler.NewAuthHandler(authService)
	noteHandler := handler.NewNoteHandler(noteService)
	checklistItemHandler := handler.NewChecklistItemHandler(checklistItemService)

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	authRouter := api.PathPrefix("/auth").Subrouter()
	authHandler.RegisterRoutes(authRouter)

	protectedRouter := api.PathPrefix("").Subrouter()
	protectedRouter.Use(handler.AuthMiddleware)

	notesRouter := protectedRouter.PathPrefix("/notes").Subrouter()
	notesRouter.HandleFunc("", noteHandler.GetNotes).Methods("GET")
	notesRouter.HandleFunc("", noteHandler.CreateNote).Methods("POST")
	notesRouter.HandleFunc("/{id:[0-9]+}", noteHandler.GetNote).Methods("GET")
	notesRouter.HandleFunc("/{id:[0-9]+}", noteHandler.UpdateNote).Methods("PUT")
	notesRouter.HandleFunc("/{id:[0-9]+}", noteHandler.DeleteNote).Methods("DELETE")

	checklistItemHandler.RegisterRoutes(notesRouter)

	// Доступен по адресу http://localhost:8080/swagger/index.html
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Printf("Сервер запускается на порту %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
