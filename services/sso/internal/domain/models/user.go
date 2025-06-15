package models

type User struct {
	ID       int64
	Email    string
	PassHash []byte
	Username string
	Role     UserRole
	Phone    string
	FarmName string
	Address  Address
}
