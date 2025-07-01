package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type MessageStore interface {
    GetSentMessages(limit, offset int) ([]Message, error)
	GetUnsentMessages(limit int) ([]Message, error)
	MarkMessageSent(id int) error
	GetDB() *sql.DB
}
 
type PostgresStore struct {
	DB *sql.DB
}

type Message struct {
	ID      int
	Content string
	Phone   string
	SentAt  sql.NullTime
}

func InitDB(dsn string) (MessageStore, error) {
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }
    db.SetMaxOpenConns(10)
    if err := db.Ping(); err != nil {
        return nil, err
    }

    return &PostgresStore{DB: db}, nil
}

func (p *PostgresStore) GetDB() *sql.DB {
	return p.DB
}

func (p *PostgresStore) GetUnsentMessages(limit int) ([]Message, error) {
	rows, err := p.DB.Query("SELECT id, content, phone FROM messages WHERE sent = FALSE LIMIT $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Content, &m.Phone); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}


func (p *PostgresStore) MarkMessageSent(id int) error {
	_, err := p.DB.Exec("UPDATE messages SET sent = TRUE, sent_at = NOW() WHERE id = $1", id)
	return err
}

func (p *PostgresStore) GetSentMessages(limit, offset int) ([]Message, error) {
	rows, err := p.DB.Query("SELECT id, content, phone, sent_at FROM messages WHERE sent = TRUE LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Content, &m.Phone, &m.SentAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

