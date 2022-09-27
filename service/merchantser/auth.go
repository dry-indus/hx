package merchantser

import (
	"fmt"
	"hx/global"
	"hx/global/context"
	"hx/mdb"
	"hx/model/merchantmod"
	"hx/service/verifyser"
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
	ErrTgExists         = fmt.Errorf("telegram already exists!")
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

	c.Infof("login success! id: %v, name: %v, category: %v", merchant.ID, merchant.Name, merchant.Category)

	return merchant, nil
}

func (AuthServer) Logout(c context.ContextB, r merchantmod.LogoutRequest) (*merchantmod.LoginResponse, error) {
	return nil, nil
}

func (AuthServer) Register(c context.ContextB, r merchantmod.RegisterRequest) (*mdb.MerchantMod, error) {
	if r.Password != r.PasswordTwo {
		return nil, ErrPwdNotMatch
	}

	err := verifyser.TgVerify.VerifyCode(c, global.RegisterSence, r.Name, r.VerifyCode)
	if err != nil {
		return nil, err
	}

	err = verifyser.TgVerify.VerifyTG(c, r.TgID, r.TgName)
	if err != nil {
		return nil, err
	}

	if count, _ := mdb.Merchant.Count(c, &mdb.MerchantTerm{TgName: &r.TgName}); count > 0 {
		return nil, ErrTgExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost) //加密处理
	if err != nil {
		c.Errorf("GenerateFromPassword failed! err: %v", err)
		return nil, ErrPwdNotMatch
	}

	mod := &mdb.MerchantMod{
		Name:      r.Name,
		Password:  string(hash),
		TgName:    r.TgName,
		TgID:      r.TgID,
		Category:  r.Category,
		CreatedAt: time.Now(),
	}

	err = mdb.Merchant.Create(c, mod)
	if err != nil {
		c.Errorf("mdb.Merchant.Create failed! err: %v", err)
		if qmgo.IsDup(err) {
			return nil, ErrAccountExists
		}
		return nil, err
	}

	c.Infof("create success! id: %v, name: %v, category: %v", mod.ID, r.Name, r.Category)

	return mod, nil
}

func (this AuthServer) FlushToken(c context.ContextB, name string) string {
	if len(name) == 0 {
		return ""
	}

	tokenKey := fmt.Sprintf(global.MERCHANT_TOEKN_KEY_FMT, name)
	token := util.UUID().String()
	global.DL_CORE_REDIS.Set(c, tokenKey, token, 8*time.Hour)

	this.flushMerchant(c, name, token)

	return token
}

func (AuthServer) flushMerchant(c context.ContextB, name, token string) {

	merchantMod, _ := mdb.Merchant.FindOneByName(c, name)
	if merchantMod == nil {
		return
	}

	merchant := context.Merchant{
		ID:       merchantMod.ID,
		Name:     merchantMod.Name,
		Category: merchantMod.Category,
		TgName:   merchantMod.TgName,
		TgID:     merchantMod.TgID,
	}

	infoKey := fmt.Sprintf(global.MERCHANT_INFO_KEY_FMT, token)
	info, _ := util.JSON.MarshalToString(merchant)
	global.DL_CORE_REDIS.Set(c, infoKey, info, 9*time.Hour)
}

func (AuthServer) RemoveToken(c context.ContextB, name string) string {
	if len(name) == 0 {
		return ""
	}

	tokenKey := fmt.Sprintf(global.MERCHANT_TOEKN_KEY_FMT, name)
	token := global.DL_CORE_REDIS.Get(c, tokenKey).Val()
	global.DL_CORE_REDIS.Del(c, tokenKey)
	return token
}

func (AuthServer) SetHoken(c context.ContextB, token string) {
	if global.Application.CloseHoken {
		return
	}

	c.Gin().Header(global.HOKEN, token)
}
