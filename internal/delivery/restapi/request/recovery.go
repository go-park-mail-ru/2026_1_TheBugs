package request

type UserEmail struct {
	Email string `json:"email"`
}

type VerifyCodeDTO struct {
	Code string `json:"code"`
}

type UpdatePwdDTO struct {
	Password string `json:"password"`
}
