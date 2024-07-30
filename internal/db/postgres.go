package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/aslbnv/messagio/internal/types"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type PostgresDB struct {
	db            *sql.DB
	dsn           string
	migrationsDir string
}

func NewPostgresDB() (*PostgresDB, error) {
	var (
		user     = viper.GetString("db.user")
		password = os.Getenv("DB_PASSWORD")
		host     = viper.GetString("db.host")
		port     = viper.GetString("db.port")
		sslmode  = viper.GetString("db.sslmode")

		migrationsDir = viper.GetString("db.migrations_dir")
		dsn           = fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=%s", user, password, host, port, sslmode)
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresDB{
		db:            db,
		dsn:           dsn,
		migrationsDir: migrationsDir,
	}, nil
}

func (p *PostgresDB) Init() error {
	m, err := migrate.New("file://"+p.migrationsDir, p.dsn)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) GetMessages() ([]*types.Message, error) {
	query := "SELECT * FROM messages"

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	messages := []*types.Message{}
	for rows.Next() {
		msg, err := scanIntoMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (p *PostgresDB) CreateMessage(msg *types.Message) error {
	query := "INSERT INTO messages (uuid, text, processed, created_at) VALUES ($1, $2, $3, $4)"

	_, err := p.db.Query(query, msg.UUID, msg.Text, msg.Processed, msg.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) GetMessageByID(id int) (*types.Message, error) {
	query := "SELECT *FROM messages WHERE id = $1"

	rows, err := p.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoMessage(rows)
	}

	return nil, fmt.Errorf("message %d not found", id)
}

func (p *PostgresDB) DeleteMessageByID(id int) error {
	query := "DELETE FROM messages WHERE id = $1"

	_, err := p.db.Query(query, id)

	return err
}

func (p *PostgresDB) MarkMessageProcessed(msg *types.Message) error {
	query := "UPDATE messages SET processed = true WHERE uuid = $1"

	_, err := p.db.Query(query, msg.UUID)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) GetProcessedMessages() (*types.ProcessedMessages, error) {
	query := "SELECT * FROM messages WHERE processed = true"

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	messages := []*types.Message{}
	for rows.Next() {
		msg, err := scanIntoMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return &types.ProcessedMessages{
		Messages: messages,
		Amount:   len(messages),
	}, nil
}

func scanIntoMessage(rows *sql.Rows) (*types.Message, error) {
	msg := new(types.Message)

	err := rows.Scan(
		&msg.ID,
		&msg.UUID,
		&msg.Text,
		&msg.Processed,
		&msg.CreatedAt,
	)

	return msg, err
}
