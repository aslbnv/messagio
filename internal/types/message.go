package types

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"id"`
	Text      string    `json:"text"`
	Processed bool      `json:"processed,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMessage(text string) *Message {
	return &Message{
		ID:        uuid.New(),
		Text:      text,
		Processed: false,
		CreatedAt: time.Now().UTC(),
	}
}

type ProcessedMessages struct {
	Amount   int        `json:"amount"`
	Messages []*Message `json:"messages"`
}

type CreateMessageRequest struct {
	Text string `json:"text"`
}

func ValidateMessageRequest(req *CreateMessageRequest) error {
	if len(req.Text) == 0 {
		return fmt.Errorf("invalid message structure")
	}

	return nil
}
