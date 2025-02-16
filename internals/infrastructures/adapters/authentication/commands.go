package authentication

import (
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
	"github.com/Pr3c10us/boilerplate/internals/domains/authentication"
	"github.com/Pr3c10us/boilerplate/packages/appError"
	"github.com/markbates/goth"
)

type RepositoryPG struct {
	db *sql.DB
}

func NewAuthenticationRepositoryPG(db *sql.DB) authentication.Repository {
	return &RepositoryPG{
		db: db,
	}
}

func (repo *RepositoryPG) CreateUser(params *authentication.AddUserParams) error {
	query, args, err := sq.Insert("users").
		Columns("email", "password", "first_name", "last_name").
		Values(params.Email, params.Password, params.FirstName, params.LastName).
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

func (repo *RepositoryPG) UpdateProfile(params *authentication.UserProfileParams) error {
	profileMap := map[string]interface{}{}

	if params.FirstName == "" {
		profileMap["first_name"] = params.FirstName
	}
	if params.LastName == "" {
		profileMap["last_name"] = params.LastName
	}
	if params.EmailVerified {
		profileMap["email_verified"] = params.EmailVerified
	}

	query, args, err := sq.Update("users").SetMap(profileMap).Where(sq.Eq{"id": params.ID}).PlaceholderFormat(sq.Dollar).ToSql()
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

func (repo *RepositoryPG) AddUserOAuth(user *goth.User) (*authentication.User, error) {
	query, args, err := sq.Insert("users").
		Columns("email", "first_name", "last_name", "email_verified").
		Values(user.Email, user.FirstName, user.LastName, true).
		Suffix(`RETURNING 
		"id",
		"email",
		"first_name",
		"last_name",
		"full_name",
		"email_verified",
		"created_at"
`).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	statement, err := repo.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var newUser authentication.User
	switch err = statement.QueryRow(args...).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.FullName,
		&newUser.EmailVerified,
		&newUser.CreatedAt,
	); {
	case errors.Is(err, sql.ErrNoRows):
		var organizationNotCreatedErr = errors.New("failed to create user")
		return nil, appError.InternalServerError(organizationNotCreatedErr)
	case err == nil:
		return &newUser, nil
	default:
		return nil, err
	}
}
