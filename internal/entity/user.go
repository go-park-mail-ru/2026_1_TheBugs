package entity

type User struct {
	ID             int
	Email          string
	Salt           string
	HashedPassword string
}
