package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aslbnv/messagio/internal/types"
)

func (s *Server) handleMessages(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodPost {
		return s.handleCreateMessage(w, r)
	}

	return fmt.Errorf("method %s not allowed", r.Method)
}

func (s *Server) handleProcessedMessages(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.GetProcessedMessages(w, r)
	}

	return fmt.Errorf("method %s not allowed", r.Method)
}

func (s *Server) handleCreateMessage(w http.ResponseWriter, r *http.Request) error {
	req := new(types.CreateMessageRequest)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	if err := types.ValidateMessageRequest(req); err != nil {
		return err
	}

	msg := types.NewMessage(req.Text)

	if err := s.db.CreateMessage(msg); err != nil {
		return err
	}

	if err := s.kafkaProducer.ProduceMessage(msg); err != nil {
		return err
	}

	if err := s.db.MarkMessageProcessed(msg); err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, msg)
}

func (s *Server) GetProcessedMessages(w http.ResponseWriter, r *http.Request) error {
	processedMessages, err := s.db.GetProcessedMessages()
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, processedMessages)
}
