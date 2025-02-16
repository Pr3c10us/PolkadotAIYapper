package embedding

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/Pr3c10us/boilerplate/internals/domains/embedding"
	"github.com/pgvector/pgvector-go"
)

type Repository struct {
	db *sql.DB
}

func NewEmbedding(db *sql.DB) embedding.Repository {
	return &Repository{db: db}
}

func (repo *Repository) AddEmbedding(embedding []float32, value string) error {
	query, args, err := sq.Insert("embeddings").
		Columns("topic", "embedding").
		Values(value, pgvector.NewVector(embedding)).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return err
	}

	var statement *sql.Stmt
	statement, err = repo.db.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(args...)
	return err
}
