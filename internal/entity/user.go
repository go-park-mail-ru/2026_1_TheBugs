package entity

type ProviderType string

var VK ProviderType = "vk"
var Yandex ProviderType = "yandex"

type User struct {
	ID             int     `db:"id"`
	Email          string  `db:"email"`
	Salt           *string `db:"salt"`
	HashedPassword *string `db:"hashed_password"`
	Provider       *string `db:"provider"`
	IsVerified     bool    `db:"is_verified"`
}

type UserDetails struct {
	ID        int     `db:"id"`
	Email     string  `db:"email"`
	FirstName string  `db:"first_name"`
	LastName  string  `db:"last_name"`
	AvatarURL *string `db:"avatar_url"`
	Phone     string  `db:"phone"`
}
