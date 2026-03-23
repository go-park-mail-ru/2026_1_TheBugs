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
}
