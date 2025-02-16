package embedding

import (
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/pgvector/pgvector-go"
)

func (repo *Repository) SimilarValuesExist(embedding []float32) (*bool, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Build query using proper parameter placeholders
	query, _, err := psql.Select("topic", "embedding <#> $1 as similarity").
		From("embeddings").
		Where("embedding <#> $1 < $2"). // Cosine similarity check
		OrderBy("embedding <#> $1").    // Order by similarity
		Limit(1).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	// Use QueryRow directly with the generated query and args
	var topic string
	var similarity float64

	err = repo.db.QueryRow(query, pgvector.NewVector(embedding), -0.7).Scan(&topic, &similarity)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		f := false
		return &f, nil
	case err == nil:
		t := true
		return &t, nil
	default:
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
}
