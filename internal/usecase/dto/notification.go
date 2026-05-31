package dto

type RecoveryNotification struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type VerificationNotification struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type AnswerNotification struct {
	Email   string `json:"email"`
	OrderID int    `json:"order_id"`
	Answer  string `json:"answer"`
}

type EmailNotification struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

type RoommateMatchNotification struct {
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PosterAlias string `json:"poster_alias"`
}

type RoommateContactsNotification struct {
	Email             string `json:"email"`
	RoommateFirstName string `json:"first_name"`
	RoommateLastName  string `json:"last_name"`
	PosterAlias       string `json:"poster_alias"`
	RoommatePhone     string `json:"phone"`
	RoommateEmail     string `json:"email_contact"`
}
