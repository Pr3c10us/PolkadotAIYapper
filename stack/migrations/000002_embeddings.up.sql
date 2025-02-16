-- Install pgvector extension (if not already installed)
CREATE EXTENSION IF NOT EXISTS vector;

-- Create the ExtraEmbeddings table
CREATE TABLE embeddings
(
    id         SERIAL PRIMARY KEY,
    topic      TEXT NOT NULL,
    embedding vector(3072),
    created_at            TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP

);
