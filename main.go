package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/epalmerini/abitudini/internal/db"
	"github.com/epalmerini/abitudini/internal/habit"
	"github.com/epalmerini/abitudini/internal/record"
	"github.com/epalmerini/abitudini/internal/streak"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	// Initialize database
	dbPath := "abitudini.db"
	database, err := db.Init(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	log.Println("Database initialized successfully")

	// Initialize slices
	// Record slice (initialize first for habit service dependency)
	recordStore := record.NewStore(database)
	recordService := record.NewService(recordStore, nil)

	// Habit slice
	habitStore := habit.NewStore(database)
	habitService := habit.NewService(habitStore, recordService)
	habitHandler := habit.NewHandler(habitService)

	// Update record service with habit service
	recordService = record.NewService(recordStore, habitService)
	recordHandler := record.NewHandler(recordService)

	// Streak slice
	streakStore := streak.NewStore(database)
	streakService := streak.NewService(streakStore)
	streakHandler := streak.NewHandler(streakService)

	// Routes
	mux := http.NewServeMux()

	// Habit API Routes
	mux.HandleFunc("POST /api/habits", habitHandler.Create)
	mux.HandleFunc("GET /api/habits", habitHandler.GetAll)
	mux.HandleFunc("GET /api/habits/{id}", habitHandler.GetByID)
	mux.HandleFunc("PUT /api/habits/{id}", habitHandler.Update)
	mux.HandleFunc("DELETE /api/habits/{id}", habitHandler.Delete)

	// Record API Routes
	mux.HandleFunc("POST /api/habits/{id}/done-today", recordHandler.MarkDoneToday)
	mux.HandleFunc("GET /api/habits/{id}/contribution", recordHandler.GetContribution)

	// Streak API Routes
	mux.HandleFunc("GET /api/habits/{id}/streak", streakHandler.GetByHabitID)

	// Static files
	staticSubFS, _ := fs.Sub(staticFiles, "static")
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticSubFS))))

	// Home page
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		habits, err := habitService.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(habit.RenderAllHabits(habits)))
	})

	// Server
	port := ":8080"
	fmt.Printf("Server running on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, mux); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
