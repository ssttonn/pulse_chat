package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
	"pulse/src/edge-ws/internal/connection"
	"pulse/src/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Cannot initialize config: %v", err)
	}

	el, _ := connection.NewEventLoop()

	go el.Start(func(conn net.Conn) {
		userID, err := connection.HandleAuthFrame(conn, cfg.JwtSecret)
		if err != nil {
			log.Printf("Auth failed for %s, closing connection: %v", conn.RemoteAddr(), err)
			conn.Close()
			return
		}

		log.Printf("User %s authenticated successfully on %s!", userID, conn.RemoteAddr())
	})

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
		log.Fatalf("Edge-WS Server crashed: %v", err)
	}
}
