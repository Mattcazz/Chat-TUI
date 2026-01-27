package user

// temporary

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

type ContactStore interface {
	GetContactByID(id int) (*Contact, error)
	CreateContact(c *Contact) error
	UpdateContact(c *Contact) error
	DeleteContact(id int) error
}

type Contact struct {
	ID       int
	UserID   int
	Username string
}
