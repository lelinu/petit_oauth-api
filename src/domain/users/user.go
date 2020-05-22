package users

type User struct {
	Id        int64 `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type UserLoginRequest struct {
	Email    string `json:"string"`
	Password string `json:"string"`
}
