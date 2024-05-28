package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type UserModel struct {
	DB *sql.DB
}

var AnonymousUser = &User{}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created-at"`
	UpdatedAt time.Time `json:"updated-at"`
	FName     string    `json:"f-name"`
	SName     string    `json:"s-name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	UserRole  string    `json:"user-role"`
	Activated bool      `json:"activated"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, err
		default:
			return false, err
		}
	}
	return true, err
}

func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO user_info (fname, sname, email, password_hash, user_role, activated)
		VALUES ($1, $2, $3, $4,$5, $6, $7)
		RETURNING id, created_at, fname, sname
	`
	args := []any{user.FName, user.SName, user.Email, user.Password.hash, user.UserRole, user.Activated}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.FName, &user.SName)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "user_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT * FROM user_info
		WHERE email = $1
	`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FName,
		&user.SName,
		&user.Email,
		&user.Password.hash,
		&user.UserRole,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) Get(id int64) (*User, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT * FROM user_info
		WHERE ID = $1
	`
	var user User
	err := m.DB.QueryRow(query, id).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.FName, &user.SName, &user.Email, &user.Password.hash, &user.UserRole, &user.Activated)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `
		UPDATE user_info
		SET fname = $1, sname = $2, email = $3, password_hash = $4, user_role = $5, activated = $6
		WHERE id = $8
		RETURNING fname, email
	`
	args := []any{user.FName, user.SName, user.Email, user.Password.hash, user.UserRole, user.Activated, user.ID}
	return m.DB.QueryRow(query, args...).Scan(&user.FName, &user.Email)
}

func (m UserModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		DELETE FROM user_info
		WHERE id = $1
	`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (m UserModel) GetAll() ([]*User, error) {
	query := `
		SELECT * FROM user_info
	`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.FName, &user.SName, &user.Email, &user.Password.hash, &user.UserRole, &user.Activated)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, ErrRecordNotFound
	}
	return users, nil
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	query := `
		SELECT user_info.id, user_info.created_at, user_info.updated_at, user_info.fname, user_info.sname, user_info.email, user_info.password_hash, user_info.user_role, user_info.activated
		FROM user_info
		INNER JOIN tokens
		ON user_info.id = tokens.user_id
		WHERE tokens.hash = $1
		AND tokens.scope = $2
		AND tokens.expiry > $3
	`

	args := []any{tokenHash[:], tokenScope, time.Now()}
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FName,
		&user.SName,
		&user.Email,
		&user.Password.hash,
		&user.UserRole,
		&user.Activated,
	)
	// token expired
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
