package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"pulse/src/edge-ws/internal/connection"
	"pulse/src/edge-ws/internal/router"
	"pulse/src/pkg/config"
	"pulse/src/pkg/kafka"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Cannot initialize config: %v", err)
	}

	producer, err := kafka.NewProducer([]string{cfg.KafkaBrokers})
	if err != nil {
		log.Fatalf("Cannot initialize kafka producer: %v", err)
	}
	defer producer.Close()

	msgRouter := router.NewRouter(producer, cfg.JwtSecret)

	el, _ := connection.NewEventLoop()

	go el.Start(msgRouter.HandleConnection)

	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := connection.UpgradeToWS(w, r)

		if err != nil {
			log.Fatalf("Cannot upgrade connection to ws: %v", err)
		}

		// Removed defer conn.Close() because kqueue event loop owns this connection now

		if err := el.Add(conn); err != nil {
			log.Printf("Failed to add to kqueue: %v", err)
			conn.Close()
			return
		}

		//nolint:gosec // address is safe
		log.Printf("Upgrade WS successfully for Client: %s", conn.RemoteAddr().String())
	})

	address := fmt.Sprintf(":%s", cfg.EdgePort)

	fmt.Printf("Edge-WS Server is running at port %s...\n", cfg.EdgePort)

	server := &http.Server{
		Addr:              address,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Printf("Edge-WS Server crashed: %v", err)
	}
}
