package models

type role int16

type User struct {
	Email    string
	Login    string
	Password string
	Role     role
	Id       int64
}
