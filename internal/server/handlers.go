package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aslbnv/messagio/internal/types"
)

func (s *Server) handleMessages(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGetMessages(w, r)
	case http.MethodPost:
		return s.handleCreateMessage(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

func (s *Server) handleMessagesByID(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGetMessageByID(w, r)
	case http.MethodDelete:
		return s.handleDeleteMessageByID(w, r)
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}
}

func (s *Server) handleProcessedMessages(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.handleGetProcessedMessages(w, r)
	}

	return fmt.Errorf("method %s not allowed", r.Method)
}

func (s *Server) handleGetMessages(w http.ResponseWriter, _ *http.Request) error {
	messages, err := s.db.GetMessages()
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, messages)
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

func (s *Server) handleGetMessageByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	msg, err := s.db.GetMessageByID(id)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, msg)
}

func (s *Server) handleDeleteMessageByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.db.DeleteMessageByID(id); err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *Server) handleGetProcessedMessages(w http.ResponseWriter, _ *http.Request) error {
	processedMessages, err := s.db.GetProcessedMessages()
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, processedMessages)
}
