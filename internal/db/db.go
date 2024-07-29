package db

import "github.com/aslbnv/messagio/internal/types"

type DB interface {
	CreateMessage(*types.Message) error
	MarkMessageProcessed(*types.Message) error
	GetProcessedMessages() (*types.ProcessedMessages, error)
}
