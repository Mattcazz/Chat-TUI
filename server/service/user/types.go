package user

type UserStore interface {
	GetUserByID(id int) (*User, error)
	CreateUser(c *User) error
	UpdateUser(c *User) error
	DeleteUser(id int) error
}

type User struct {
	ID       int
	Username string
}
