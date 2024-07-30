# Messagio
A simple RESTful API for handling messages on Golang

## Requirements
- Docker
- Docker Compose

## Run
Just exec:
```sh
make up
```

If you want to specify your own free port values for containers, create `.env` file as shown below:
```env
APP_PORT=3000
DB_PORT=5436
KAFKA_PORT=9092
```

## Structure
```sh
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
- `GET` : Get all messages
- `POST` : Create a new message
```json
{
    "text":"hello"
}
```

#### /messages/id
- `GET` : Get selected message
- `DELETE` : Delete selected message

#### /messages/processed
- `GET` : Get only processed messages