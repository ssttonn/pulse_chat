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

// ReadFrame reads a WS frame, unmasks it, and returns the raw payload.
// It returns a 'release' function that the caller MUST defer to prevent memory leaks!
func ReadFrame(conn net.Conn) ([]byte, func(), error) {
	header, err := ws.ReadHeader(conn)
	if err != nil {
		return nil, nil, fmt.Errorf("read header: %v", err)
	}

	// Cap payload size at 8KB to prevent OOM
	if header.Length > 8192 {
		return nil, nil, fmt.Errorf("payload too large: %d", header.Length)
	}

	bufPtr := pool.GetBuffer()
	if bufPtr == nil {
		return nil, nil, fmt.Errorf("failed to get buffer")
	}

	release := func() {
		pool.PutBuffer(bufPtr)
	}

	payload := (*bufPtr)[:header.Length]

	if _, err := io.ReadFull(conn, payload); err != nil {
		release()
		return nil, nil, fmt.Errorf("read payload: %v", err)
	}

	if header.Masked {
		ws.Cipher(payload, header.Mask, 0)
	}

	return payload, release, nil
}

func HandleAuthFrame(conn net.Conn, secret string) (string, error) {
	payload, release, err := ReadFrame(conn)
	if err != nil {
		return "", err
	}
	defer release()

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
