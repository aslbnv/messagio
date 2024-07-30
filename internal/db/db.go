package db

import "github.com/aslbnv/messagio/internal/types"

type DB interface {
	GetMessages() ([]*types.Message, error)
	CreateMessage(*types.Message) error
	GetMessageByID(int) (*types.Message, error)
	DeleteMessageByID(int) error
	GetProcessedMessages() (*types.ProcessedMessages, error)
	MarkMessageProcessed(*types.Message) error
}
