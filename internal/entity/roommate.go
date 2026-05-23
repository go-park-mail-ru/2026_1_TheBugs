package entity

type RoommateUser struct {
	FirstName   string
	LastName    string
	AvatarURL   *string
	Gender      string
	Birthday    string
	Description *string
}

type RoommateTag struct {
	Name  string
	Alias string
}

type RoommateForm struct {
	Gender      string
	Birthday    string
	Description string
}
