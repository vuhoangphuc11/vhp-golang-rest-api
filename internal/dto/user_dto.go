package dto

type UserDto struct {
	Username        string
	Email           string
	FirstName       string
	LastName        string
	Password        string
	ConfirmPassword string
	Age             int
	Gender          bool
	Phone           string
	IsActive        bool
	Role            string
}
