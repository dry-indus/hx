package merchantser

import (
	"hx/global/context"
	"hx/model/merchantmod"
)

var (
	Auth AuthServer
)

type AuthServer struct {
}

func (AuthServer) Login(c context.ContextB, r merchantmod.LoginRequest) (*merchantmod.RegisterResponse, error) {
	return nil, nil
}

func (AuthServer) Logout(c context.ContextB, r merchantmod.LogoutRequest) (*merchantmod.LoginResponse, error) {
	return nil, nil
}

func (AuthServer) Register(c context.ContextB, r merchantmod.RegisterRequest) (*merchantmod.RegisterResponse, error) {
	return nil, nil
}
