package entity

type (
	ReqID     struct{}
	LoggerCtx struct{}
	UserID    struct{}
	Metre     uint
)

type ServiceType string

const (
	GatewayService ServiceType = "gateaway_service"
	AuthService    ServiceType = "auth_service"
	PosterService  ServiceType = "poster_service"
	UserService    ServiceType = "user_service"
)

var (
	RecoveryRoute                     = "recovery_code"
	VerificationRoute                 = "verification_code"
	AnswerRoute                       = "answer_notification"
	EmailRoute                        = "email_notification"
	RoommateMatchRoute                = "roommate_match"
	RoommateContactsForRequesterRoute = "roommate_contacts_requester"
	RoommateContactsForAcceptedRoute  = "roommate_contacts_accepted"
)
