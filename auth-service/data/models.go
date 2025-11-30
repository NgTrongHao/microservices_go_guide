package data

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 10

//func New(dbPool *sql.DB) Models {
//	db = dbPool
//
//	return Models{
//		User: User{},
//	}
//}
//
//type Models struct {
//	User User
//}

type UserRepository struct {
	Conn *sql.DB
}

func New(dbPool *sql.DB) *UserRepository {
	return &UserRepository{Conn: dbPool}
}

type User struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	FirstName string     `json:"first_name,omitempty"`
	LastName  string     `json:"last_name,omitempty"`
	Active    bool       `json:"active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (r *UserRepository) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, password, first_name, last_name, active, created_at, updated_at, deleted_at FROM users`

	rows, err := r.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}(rows)

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.FirstName,
			&user.LastName,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *UserRepository) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `SELECT id, email, password, first_name, last_name, active, created_at, updated_at, deleted_at 
			FROM users 
			WHERE email = $1`
	row := r.Conn.QueryRowContext(ctx, query, email)
	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `SELECT id, email, password, first_name, last_name, active, created_at, updated_at, deleted_at 
			FROM users 
			WHERE id = $1`
	row := r.Conn.QueryRowContext(ctx, query, id)
	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(u *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `UPDATE users 
			SET email = $1, first_name = $2, last_name = $3, active = $4, updated_at = $5
			WHERE id = $6`
	_, err := r.Conn.ExecContext(ctx, query,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Active,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Delete(u *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `UPDATE users 
			SET deleted_at = ?, active = false 
			WHERE id = ?`
	_, err := r.Conn.ExecContext(ctx, query, time.Now(), u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) DeleteById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `UPDATE users 
			SET deleted_at = ?, active = false 
			WHERE id = ?`
	_, err := r.Conn.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Insert(u *User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var id int
	stmt := `INSERT INTO users 
    		(email, password, first_name, last_name, active) 
			VALUES (?, ?, ?, ?, ?) 
			RETURNING id`

	err = r.Conn.QueryRowContext(ctx, stmt, u.Email, hashedPassword, u.FirstName, u.LastName, u.Active).Scan(&id)
	if err != nil {
		return 0, err
	}
	u.ID = id
	return id, nil
}

func (r *UserRepository) ValidatePassword(password string, u *User) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *UserRepository) ResetPassword(password string, u *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `UPDATE users 
			SET password = ?, updated_at = ?
			WHERE id = ?`
	_, err = r.Conn.ExecContext(ctx, query,
		hashedPassword,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
