package api

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/reddtsai/goAPI/pkg/blockaction/storage"
)

type UserClaims struct {
	ID      int64  `json:"id"`
	Account string `json:"account"`
	jwt.RegisteredClaims
}

type SignupReq struct {
	Account  string `json:"account" binding:"required,min=6,max=12"`   // 帳號
	Password string `json:"password" binding:"required,min=8,max=16"`  // 密碼
	UserName string `json:"user_name" binding:"required,min=2,max=20"` // 名稱
}

func (c *SignupReq) ToEntity() (entity storage.UserTable, err error) {
	id := _snowNode.Generate().Int64()
	entity = storage.UserTable{
		ID:        id,
		Account:   c.Account,
		Name:      c.UserName,
		CreatedAt: time.Now().UnixMilli(),
		Creator:   id,
		UpdatedAt: time.Now().UnixMilli(),
		Updater:   id,
	}
	entity.Secret, err = signaturePwd(c.Password)

	return
}

type SignupResp struct {
	Account string `json:"account"`
}

type SigninReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SigninResp struct {
	Token        string `json:"token"`         // token
	RefreshToken string `json:"refresh_token"` // refresh token
}

type GetPersonalInfoResp struct {
	ID       string `json:"id"`
	Account  string `json:"account"`
	UserName string `json:"user_name"`
}
