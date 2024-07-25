package worker

import (
	"context"
	"log"
	"pulse/src/message-router/internal/db"
	"pulse/src/pkg/models"
	"time"
)

type Batcher struct {
	messageChan chan models.ChatPayload
	batchSize   int
	repo        *db.Repository
}

func NewBatcher(maxBufferSize int, batchSize int, repo *db.Repository) *Batcher {
	return &Batcher{
		messageChan: make(chan models.ChatPayload, maxBufferSize),
		batchSize:   batchSize,
		repo:        repo,
	}
}

func (b *Batcher) Add(msg models.ChatPayload) {
	b.messageChan <- msg
}

func (b *Batcher) Start(ctx context.Context) {
	go func() {
		batch := make([]models.ChatPayload, 0, b.batchSize)

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case newMsg := <-b.messageChan:
				batch = append(batch, newMsg)
				if len(batch) >= b.batchSize {
					batch = batch[:0]
				}
			case <-ticker.C:
				if len(batch) > 0 {
					log.Printf("Flushing batch of size %d", len(batch))
					batch = batch[:0]
				}
			}
		}
	}()
}
