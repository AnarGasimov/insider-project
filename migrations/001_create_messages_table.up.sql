CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    content VARCHAR(160) NOT NULL,
    phone VARCHAR(15) NOT NULL,
    sent BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP
);
CREATE INDEX idx_sent ON messages(sent);
