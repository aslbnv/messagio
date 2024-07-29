package main

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/aslbnv/messagio/internal/db"
	"github.com/aslbnv/messagio/internal/kafka"
	"github.com/aslbnv/messagio/internal/server"
)

var (
	configPath = "configs"
	configName = "config"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %s", err)
	}

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err)
	}
	
	db, err := db.NewPostgresDB()
	if err != nil {
		log.Fatalf("error connecting to database: %s", err)
	}

	if err := db.Init(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("error migrating database: %s", err)
	}

	kafkaProducer, err := kafka.NewKafkaProducer()
	if err != nil {
		log.Fatalf("error connecting to kafka: %s", err)
	}

	server := server.NewServer(
		viper.GetString("server.port"),
		db,
		kafkaProducer,
	)

	log.Printf("server running on port :%s\n", viper.GetString("server.port"))
	if err := server.Start(); err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}

func initConfig() error {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)

	return viper.ReadInConfig()
}
