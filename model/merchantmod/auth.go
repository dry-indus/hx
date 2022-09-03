package merchantmod

type LoginRequest struct {
}

type LoginResponse struct {
}

type LogoutRequest struct {
}

type LogoutResponse struct {
}

type RegisterRequest struct {
	Name     string
	Password string
	Code     string
	Class    string
}

type RegisterResponse struct {
}
