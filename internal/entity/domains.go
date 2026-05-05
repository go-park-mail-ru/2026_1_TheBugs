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
