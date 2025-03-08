package app_interfaces

type UserManager interface {
	CheckTgUserAuth(id int) (bool, error)
	Auth(id int, email string, password []byte) error
	Logout(id int) error
}
