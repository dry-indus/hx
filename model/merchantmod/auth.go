package merchantmod

type LoginRequest struct {
	Name     string
	Password string
	Code     string
}

type LoginResponse struct {
	Name     string
	Telegram string
	Class    int
}

type LogoutRequest struct {
}

type LogoutResponse struct {
}

type RegisterRequest struct {
	Name        string
	Password    string
	PasswordTwo string
	Code        string
	Telegram    string
	Category    int
}

type RegisterResponse struct {
	Name     string
	Telegram string
	Category int
}
