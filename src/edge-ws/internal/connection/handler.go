package connection

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"pulse/src/edge-ws/internal/models"
	"pulse/src/edge-ws/internal/pool"
	"pulse/src/pkg/auth"

	"github.com/gobwas/ws"
)

func HandleAuthFrame(conn net.Conn, secret string) (string, error) {
	header, err := ws.ReadHeader(conn)
	if err != nil {
		return "nil", fmt.Errorf("failed to read ws header: %v", err)
	}

	if header.Length > 4096 {
		return "", fmt.Errorf("auth payload too large: %d", header.Length)
	}

	bufPtr := pool.GetBuffer()

	if bufPtr == nil {
		return "", fmt.Errorf("failed to get buffer from pool")
	}

	defer pool.PutBuffer(bufPtr)

	payload := (*bufPtr)[:header.Length]

	if _, err := io.ReadFull(conn, payload); err != nil {
		return "", fmt.Errorf("failed to read payload: %v", err)
	}

	if header.Masked {
		ws.Cipher(payload, header.Mask, 0)
	}

	var inbound models.InboundMessage

	if err := json.Unmarshal(payload, &inbound); err != nil {
		return "", fmt.Errorf("invalid json format: %v", err)
	}

	if inbound.Type != models.MessageTypeAuth {
		return "", fmt.Errorf("expected auth message, got %s", inbound.Type)
	}

	var authPayload models.AuthPayload

	if err := json.Unmarshal(inbound.Payload, &authPayload); err != nil {
		return "", fmt.Errorf("invalid auth payload: %v", err)
	}

	userID, err := auth.ValidateToken(authPayload.Token, secret)
	if err != nil {
		return "", fmt.Errorf("invalid token: %v", err)
	}

	return userID, nil
}
