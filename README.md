# Messagio
A simple RESTful API for handling messages on Golang

## Requirements
- Docker
- Docker Compose v2.x.x

## Run
Create ```.env``` file in root directory of project with values of **free** ports on your machine and PG password
```env
APP_PORT=3000 # free port to run Go app
DB_PORT=5432 # free port to run Postgres
KAFKA_PORT=9092 # free port to run Kafka

POSTGRES_PASSWORD=password
```
Then
```sh
make up
```

## Structure
```
├── cmd
│   └── main.go
├── configs
│   └── config.yml
├── internal
│   ├── db
│   │   ├── db.go
│   │   └── postgres.go
│   ├── kafka
│   │   └── producer.go
│   ├── server
│   │   ├── handlers.go
│   │   ├── server.go
│   │   └── utils.go
│   └── types
│       ├── api.go
│       └── message.go
└── migrations
    ├── 000001_create_message_table.down.sql
    └── 000001_create_message_table.up.sql
```

## API
#### /messages
- `POST` : Create a new message
```json
{
    "text":"hello"
}
```

#### /messages/processed
- `GET` : Get processed messages