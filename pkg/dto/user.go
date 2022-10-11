package dto

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Password string `json:"password"`
	Age      int    `json:"age"`
	Gender   bool   `json:"gender"`
	Phone    string `json:"phone"`
	IsActive bool   `json:"active"`
	Role     string `json:"role"`
}
