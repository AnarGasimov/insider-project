package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

 
type Message struct {
	ID      int
	Content string
	Phone   string
	SentAt  sql.NullTime
}

func InitDB(dsn string) error {
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	DB.SetMaxOpenConns(10) 
	return DB.Ping()
}


func GetUnsentMessages(limit int) ([]Message, error) {
	rows, err := DB.Query("SELECT id, content, phone FROM messages WHERE sent = FALSE LIMIT $1", limit)
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


func MarkMessageSent(id int) error {
	_, err := DB.Exec("UPDATE messages SET sent = TRUE, sent_at = NOW() WHERE id = $1", id)
	return err
}

func GetSentMessages(limit, offset int) ([]Message, error) {
	rows, err := DB.Query("SELECT id, content, phone, sent_at FROM messages WHERE sent = TRUE LIMIT $1 OFFSET $2", limit, offset)
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

