package entity

type User struct {
	ID             int    `db:"id"`
	Email          string `db:"email"`
	Salt           string `db:"salt"`
	HashedPassword string `db:"hashed_password"`
}
