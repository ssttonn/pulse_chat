package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"pulse/src/api-core/internal/adapter/db"
	apphttp "pulse/src/api-core/internal/adapter/http"
	"pulse/src/api-core/internal/usecase"
	"pulse/src/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Không thể khởi tạo config: %v", err)
	}

	// 1. Setup Database Connection (Postgres)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost, cfg.PostgresUser, cfg.PostgresPass, cfg.PostgresDB)

	database, err := db.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("Cannot connect to db: %v", err)
	}
	defer database.Close()

	// 2. Setup Repositories (Data Layer)
	userRepo := db.NewUserRepository(database)

	// 3. Setup UseCases (Business Layer)
	userUseCase := usecase.NewUserUseCase(userRepo)

	// 4. Setup Handlers (Transport Layer)
	userHandler := apphttp.NewUserHandler(userUseCase)

	// 5. Inject Handlers into Router
	router := apphttp.NewRouter(userHandler)

	address := fmt.Sprintf(":%s", cfg.APIPort)
	fmt.Printf("API-Core Server is running on port %s...\n", cfg.APIPort)

	server := &http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: 3 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("Server crashed: %v", err)
	}
}
