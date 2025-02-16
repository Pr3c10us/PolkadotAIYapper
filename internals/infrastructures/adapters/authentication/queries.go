package authentication

import (
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/packages/appError"
	"github.com/google/uuid"
)

func (repo *RepositoryPG) GetUserDetails(params *authentication.GetUserParams) (*authentication.User, error) {
	query, args, err := sq.Select(
		"id",
		"COALESCE(email, '') AS email",
		"password",
		"COALESCE(first_name, '') AS first_name",
		"COALESCE(last_name, '') AS last_name",
		"COALESCE(full_name, '') AS full_name",
		"created_at",
	).From("users").Where(sq.Or{sq.Eq{"email": params.Email}, sq.Eq{"id": params.ID}}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var statement *sql.Stmt
	statement, err = repo.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	user := authentication.User{}
	var id string
	switch err = statement.QueryRow(args...).Scan(
		&id,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.FullName,
		&user.CreatedAt,
	); {
	case errors.Is(err, sql.ErrNoRows):
		var userNotFoundErr = errors.New("user with does not exit")
		return nil, appError.NotFound(userNotFoundErr)
	case err == nil:
		user.ID, err = uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		return &user, nil
	default:
		return nil, err
	}
}
