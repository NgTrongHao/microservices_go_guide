package data

type Repository interface {
	GetAll() ([]*User, error)
	GetByEmail(email string) (*User, error)
	GetByID(id int) (*User, error)
	Update(user *User) error
	Delete(user *User) error
	DeleteById(id int) error
	Insert(user *User) (int, error)
	ValidatePassword(password string, user *User) (bool, error)
	ResetPassword(password string, user *User) error
}
