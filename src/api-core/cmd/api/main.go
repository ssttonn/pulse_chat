package main

import (
	"fmt"
	"log"
	"net/http"
	apphttp "pulse/src/api-core/internal/adapter/http"
	"pulse/src/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Không thể khởi tạo config: %v", err)
	}
	router := apphttp.NewRouter()

	address := fmt.Sprintf(":%s", cfg.APIPort)
	fmt.Printf("Server is running on port %s...", cfg.APIPort)

	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatalf("Server bị sập: %v", err)
	}
}
