package functions

import (
	"CatsSocial/configs"
	"CatsSocial/db/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	config configs.Config
	dbPool *pgxpool.Pool
}

func NewUser(dbPool *pgxpool.Pool, config configs.Config) *User {
	return &User{
		dbPool: dbPool,
		config: config,
	}
}

func (u *User) Register(ctx context.Context, usr models.User) (models.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return models.User{}, err
	}
	defer conn.Release()

	// Hash the password before storing it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(usr.Password), u.config.BcryptSalt)
	if err != nil {
		return models.User{}, err
	}

	var existingId string

	err = conn.QueryRow(ctx, `SELECT id FROM users WHERE email = $1`, usr.Email).Scan(&existingId)
	if existingId != "" {
		return models.User{}, errors.New("EXISTING_EMAIL")
	}

	sql := `
		INSERT INTO users (email, name, password) VALUES ($1, $2, $3)
	`

	_, err = conn.Exec(ctx, sql, usr.Email, usr.Name, string(hashedPassword))

	var result models.User

	err = conn.QueryRow(ctx, `SELECT id, email, name FROM users WHERE email = $1`, usr.Email).Scan(&result.Id, &result.Email, &result.Name)

	if err != nil {
		return models.User{}, err
	}

	return models.User{
		Id:    result.Id,
		Email: result.Email,
		Name:  result.Name,
	}, nil
}

func (u *User) Login(ctx context.Context, email, password string) (models.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return models.User{}, err
	}
	defer conn.Release()

	var result models.User

	err = conn.QueryRow(ctx, `SELECT id, email, name, password FROM users WHERE email = $1`, email).Scan(
		&result.Id, &result.Email, &result.Name, &result.Password,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, errors.New("USER_NOT_FOUND")
	}
	if err != nil {
		return result, err
	}

	// Compare the provided password with the hashed password from the database
	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)); err != nil {
		return result, errors.New("INVALID_PASSWORD")
	}

	return result, nil
}

func (u *User) GetUserById(ctx context.Context, userID string) (models.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		return models.User{}, err
	}
	defer conn.Release()

	var result models.User

	err = conn.QueryRow(ctx, `SELECT id, email, name FROM users WHERE id = $1`, userID).Scan(&result.Id, &result.Email, &result.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, ErrNoRow
	}
	if err != nil {
		return result, err
	}

	return result, nil
}
