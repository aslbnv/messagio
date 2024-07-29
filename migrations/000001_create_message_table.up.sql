CREATE TABLE messages (
    id UUID PRIMARY KEY,
    text TEXT NOT NULL,
    processed BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);
