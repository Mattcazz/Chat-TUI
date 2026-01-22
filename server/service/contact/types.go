package contact

type ContactStore interface {
	GetContactByID(id int) (*Contact, error)
	CreateContact(c *Contact) error
	UpdateContact(c *Contact) error
	DeleteContact(id int) error
}

type Contact struct {
	ID       int
	Username string
}
