package mango

// Identity is the interface used to hold information of the user
// making the request. Implementations of Identity are should to created
// and populated in a pre-hook.
type Identity interface {
	UserID() string
	Email() string
	Fullname() string
	Organization() string
}

// BasicIdentity is a basic implementation of the Identity interface.
type BasicIdentity struct {
	Username string
}

// UserID returns the ID of the request user
func (i BasicIdentity) UserID() string { return i.Username }

// Email returns the email address of the request user
func (i BasicIdentity) Email() string { return "" }

// Fullname returns the fullname of the request user
func (i BasicIdentity) Fullname() string { return "" }

// Organization returns the organization of the request user
func (i BasicIdentity) Organization() string { return "" }
