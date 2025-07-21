package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gocbt/internal/api"
	"gocbt/internal/auth"
	"gocbt/internal/config"
	"gocbt/internal/database"
	"gocbt/internal/middleware"
	"gocbt/internal/services"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations("migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := database.NewUserRepository(db)
	testRepo := database.NewTestRepository(db)
	questionRepo := database.NewQuestionRepository(db)
	sessionRepo := database.NewTestSessionRepository(db)
	answerRepo := database.NewUserAnswerRepository(db)
	resultRepo := database.NewTestResultRepository(db)

	// Initialize services
	passwordManager := auth.NewPasswordManager()
	userService := services.NewUserService(userRepo, passwordManager)
	testService := services.NewTestService(testRepo)
	questionService := services.NewQuestionService(questionRepo)
	resultService := services.NewTestResultService(resultRepo, sessionRepo, answerRepo, testRepo, questionRepo)
	sessionService := services.NewTestSessionService(sessionRepo, answerRepo, testRepo, questionRepo, resultService)

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(&cfg.JWT)

	// Initialize middleware
	authMiddleware := auth.NewMiddleware(jwtManager)

	// Initialize handlers
	authHandler := api.NewAuthHandler(userService, jwtManager)
	testHandler := api.NewTestHandler(testService, questionService)
	questionHandler := api.NewQuestionHandler(questionService)
	sessionHandler := api.NewSessionHandler(sessionService)
	resultHandler := api.NewResultHandler(resultService)

	// Setup routes
	router := setupRoutes(authHandler, testHandler, questionHandler, sessionHandler, resultHandler, authMiddleware)

	// Create rate limiter (100 requests per minute per IP)
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	// Apply security middleware
	secureRouter := middleware.SecurityHeaders(router)
	secureRouter = rateLimiter.RateLimit(secureRouter)
	secureRouter = middleware.RequestSizeLimit(10 * 1024 * 1024)(secureRouter) // 10MB limit
	secureRouter = middleware.ValidateContentType("application/json")(secureRouter)

	// Setup CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins(cfg.App.CORSOrigins),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(secureRouter)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      corsHandler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// setupRoutes configures the application routes
func setupRoutes(authHandler *api.AuthHandler, testHandler *api.TestHandler, questionHandler *api.QuestionHandler, sessionHandler *api.SessionHandler, resultHandler *api.ResultHandler, authMiddleware *auth.Middleware) *mux.Router {
	router := mux.NewRouter()

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	// API routes
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	// Auth routes (public)
	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", authHandler.Register).Methods("POST")
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Protected auth routes
	protectedAuthRouter := apiRouter.PathPrefix("/auth").Subrouter()
	protectedAuthRouter.Use(authMiddleware.Authenticate)
	protectedAuthRouter.HandleFunc("/profile", authHandler.Profile).Methods("GET")
	protectedAuthRouter.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")

	// Test routes (protected)
	testRouter := apiRouter.PathPrefix("/tests").Subrouter()
	testRouter.Use(authMiddleware.Authenticate)
	testRouter.HandleFunc("", testHandler.CreateTest).Methods("POST")
	testRouter.HandleFunc("", testHandler.ListTests).Methods("GET")
	testRouter.HandleFunc("/available", testHandler.GetAvailableTests).Methods("GET")
	testRouter.HandleFunc("/{id:[0-9]+}", testHandler.GetTest).Methods("GET")
	testRouter.HandleFunc("/{id:[0-9]+}", testHandler.UpdateTest).Methods("PUT")
	testRouter.HandleFunc("/{id:[0-9]+}", testHandler.DeleteTest).Methods("DELETE")
	testRouter.HandleFunc("/{id:[0-9]+}/questions", testHandler.GetTestQuestions).Methods("GET")

	// Question routes (protected)
	questionRouter := apiRouter.PathPrefix("/questions").Subrouter()
	questionRouter.Use(authMiddleware.Authenticate)
	questionRouter.HandleFunc("", questionHandler.CreateQuestion).Methods("POST")
	questionRouter.HandleFunc("/{id:[0-9]+}", questionHandler.GetQuestion).Methods("GET")
	questionRouter.HandleFunc("/{id:[0-9]+}", questionHandler.UpdateQuestion).Methods("PUT")
	questionRouter.HandleFunc("/{id:[0-9]+}", questionHandler.DeleteQuestion).Methods("DELETE")
	questionRouter.HandleFunc("/{id:[0-9]+}/options", questionHandler.AddOption).Methods("POST")
	questionRouter.HandleFunc("/{id:[0-9]+}/options/{optionId:[0-9]+}", questionHandler.UpdateOption).Methods("PUT")
	questionRouter.HandleFunc("/{id:[0-9]+}/options/{optionId:[0-9]+}", questionHandler.DeleteOption).Methods("DELETE")
	questionRouter.HandleFunc("/{id:[0-9]+}/answers", questionHandler.AddCorrectAnswer).Methods("POST")

	// Session routes (protected)
	sessionRouter := apiRouter.PathPrefix("/sessions").Subrouter()
	sessionRouter.Use(authMiddleware.Authenticate)
	sessionRouter.HandleFunc("/start", sessionHandler.StartSession).Methods("POST")
	sessionRouter.HandleFunc("/my", sessionHandler.GetUserSessions).Methods("GET")
	sessionRouter.HandleFunc("/{token}", sessionHandler.GetSession).Methods("GET")
	sessionRouter.HandleFunc("/{token}/answers", sessionHandler.SubmitAnswer).Methods("POST")
	sessionRouter.HandleFunc("/{token}/answers", sessionHandler.GetSessionAnswers).Methods("GET")
	sessionRouter.HandleFunc("/{token}/submit", sessionHandler.SubmitSession).Methods("POST")
	sessionRouter.HandleFunc("/{token}/progress", sessionHandler.UpdateProgress).Methods("PUT")

	// Result routes (protected)
	resultRouter := apiRouter.PathPrefix("/results").Subrouter()
	resultRouter.Use(authMiddleware.Authenticate)
	resultRouter.HandleFunc("/my", resultHandler.GetUserResults).Methods("GET")
	resultRouter.HandleFunc("/{id:[0-9]+}", resultHandler.GetResult).Methods("GET")
	resultRouter.HandleFunc("/session/{sessionId:[0-9]+}", resultHandler.GetResultBySession).Methods("GET")
	resultRouter.HandleFunc("/session/{sessionId:[0-9]+}/calculate", resultHandler.CalculateResult).Methods("POST")
	resultRouter.HandleFunc("/test/{id:[0-9]+}", resultHandler.GetTestResults).Methods("GET")
	resultRouter.HandleFunc("/test/{id:[0-9]+}/statistics", resultHandler.GetTestStatistics).Methods("GET")

	return router
}
