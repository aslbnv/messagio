CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL,
    text TEXT NOT NULL,
    processed BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);
