package user

//User implements user struct
type User struct {
	Username string
	Password string
}

//New returns new username with given password
func New(username string, password string) User {
	return User{username, password}
}
