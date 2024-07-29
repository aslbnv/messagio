package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/aslbnv/messagio/internal/db"
	"github.com/aslbnv/messagio/internal/kafka"
)

type Server struct {
	listenAddr    string
	db            db.DB
	kafkaProducer *kafka.KafkaProducer
}

func NewServer(listenAddr string, db db.DB, kafkaProducer *kafka.KafkaProducer) *Server {
	return &Server{
		listenAddr:    ":" + listenAddr,
		db:            db,
		kafkaProducer: kafkaProducer,
	}
}

func (s *Server) Start() error {
	router := mux.NewRouter()

	router.HandleFunc("/messages", makeHandler(s.handleMessages))
	router.HandleFunc("/messages/processed", makeHandler(s.handleProcessedMessages))

	return http.ListenAndServe(s.listenAddr, router)
}
