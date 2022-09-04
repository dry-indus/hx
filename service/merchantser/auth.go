package merchantser

import (
	"fmt"
	"hx/global"
	"hx/global/context"
	"hx/mdb"
	"hx/model/merchantmod"
	"hx/util"
	"time"

	"github.com/qiniu/qmgo"
	"golang.org/x/crypto/bcrypt"
)

var (
	Auth AuthServer
)

type AuthServer struct {
}

var (
	ErrAccountExists    = fmt.Errorf("account already exists!")
	ErrAccountNotExists = fmt.Errorf("account not exists!")
	ErrCreate           = fmt.Errorf("create account failed!")
	ErrPwdNotMatch      = fmt.Errorf("password does not match!")
)

func (AuthServer) Login(c context.ContextB, r merchantmod.LoginRequest) (*mdb.MerchantMod, error) {
	merchant, err := mdb.Merchant.FindOneByName(c, r.Name)
	if err != nil {
		if qmgo.IsErrNoDocuments(err) {
			return nil, ErrAccountNotExists
		}
		c.Errorf("mdb.Merchant.Create faield! err: %v", err)
		return nil, ErrAccountNotExists
	}

	err = bcrypt.CompareHashAndPassword([]byte(merchant.Password), []byte(r.Password)) //验证（对比）
	if err != nil {
		return nil, ErrPwdNotMatch
	}

	c.Infof("login success! id: %v, name: %v, class: %v", merchant.ID, merchant.Name, merchant.Class)

	return merchant, nil
}

func (AuthServer) Logout(c context.ContextB, r merchantmod.LogoutRequest) (*merchantmod.LoginResponse, error) {
	return nil, nil
}

func (AuthServer) Register(c context.ContextB, r merchantmod.RegisterRequest) error {
	if r.Password != r.PasswordTwo {
		return ErrPwdNotMatch
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		c.Errorf("GenerateFromPassword failed! err: %v", err)
		return ErrPwdNotMatch
	}
	encodePWD := string(hash)

	mod := mdb.MerchantMod{
		Name:      r.Name,
		Password:  encodePWD,
		Telegram:  r.Telegram,
		Class:     r.Class,
		CreatedAt: time.Now(),
	}

	id, err := mdb.Merchant.Create(c, mod)
	if err != nil {
		if qmgo.IsDup(err) {
			return ErrAccountExists
		}
		c.Errorf("mdb.Merchant.Create faield! err: %v", err)
		return ErrCreate
	}

	c.Infof("create success! id: %v, name: %v, class: %v", id, r.Name, r.Class)

	return nil
}

func (AuthServer) FlushToken(c context.ContextB, merchant *mdb.MerchantMod) string {
	tokenKey := fmt.Sprintf(global.MERCHANT_TOEKN_KEY_FMT, merchant.Name)
	token := util.UUID().String()
	global.DL_CORE_REDIS.Set(c, tokenKey, token, 8*time.Hour)
	return token
}
