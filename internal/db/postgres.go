package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/aslbnv/messagio/internal/types"
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

func (p *PostgresDB) CreateMessage(msg *types.Message) error {
	query := "INSERT INTO messages (id, text, processed, created_at) VALUES ($1, $2, $3, $4)"

	_, err := p.db.Query(query, msg.ID, msg.Text, msg.Processed, msg.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresDB) MarkMessageProcessed(msg *types.Message) error {
	query := "UPDATE messages SET processed = true WHERE id = $1"

	_, err := p.db.Query(query, msg.ID)
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
		&msg.Text,
		&msg.Processed,
		&msg.CreatedAt,
	)

	return msg, err
}
