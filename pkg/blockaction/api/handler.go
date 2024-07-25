package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type BlockActionApi struct {
	*gin.Engine
	opts BlockActionApiOptions
}

func (b *BlockActionApi) Health(c *gin.Context) {
	c.Status(http.StatusOK)
}

// @Summary 會員註冊
// @Description 會員註冊
// @Tags BlockAction
// @Accept json
// @Produce json
// @Param Request body SignupReq true "raw"
// @Success 200 {object} BaseResponse{result=SignupResp} "ok"
// @Failure 400 {object} BaseResponse "bad request"
// @Failure 403 {object} BaseResponse "forbidden"
// @Failure 409 {object} BaseResponse "conflict"
// @Failure 500 {object} BaseResponse "server error"
// @Router /v1/signup [post]
func (b *BlockActionApi) Signup(c *gin.Context) {
	req := SignupReq{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	entity, err := req.ToEntity()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	exist, err := b.opts.storage.IsExistUserAccount(req.Account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exist {
		c.JSON(http.StatusConflict, gin.H{"error": "account already exist"})
		return
	}
	err = b.opts.storage.CreateUser(entity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": "0",
		"result": SignupResp{
			Account: entity.Account,
		},
	})
}

// @Summary 會員登入
// @Description 會員登入
// @Tags BlockAction
// @Accept json
// @Produce json
// @Param Request body SigninReq true "raw"
// @Success 200 {object} BaseResponse{result=SigninResp} "ok"
// @Failure 400 {object} BaseResponse{error=string} "bad request"
// @Failure 403 {object} BaseResponse{error=string} "forbidden"
// @Failure 500 {object} BaseResponse{error=string} "server error"
// @Router /v1/signin [post]
func (b *BlockActionApi) Signin(c *gin.Context) {
	req := &SigninReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u, err := b.opts.storage.GetUserByAccount(req.Account)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if u.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "account or password incorrect"})
		return
	}
	ok, err := validatePwd(req.Password, u.Secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "account or password incorrect"})
		return
	}

	expiresTime := time.Duration(TOKEN_EXPIRE_TIME)
	token, err := genToken(&UserClaims{
		ID:      u.ID,
		Account: u.Account,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "blockaction",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresTime * time.Second)),
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	reToken, err := genToken(&UserClaims{
		ID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "blockaction",
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": &SigninResp{
			Token:        token,
			RefreshToken: reToken,
		},
	})
}

// @Summary 會員資訊
// @Description 會員資訊
// @Tags BlockAction
// @Accept json
// @Produce json
// @Success 200 {object} BaseResponse{result=GetPersonalInfoResp} "ok"
// @Failure 401 {object} BaseResponse "unauthorized"
// @Failure 403 {object} BaseResponse "forbidden"
// @Failure 500 {object} BaseResponse "server error"
// @Router /v1/user/personal-info [get]
func (b *BlockActionApi) GetPersonalInfo(c *gin.Context) {
	userID := c.GetInt64(CTX_USER_ID)
	u, err := b.opts.storage.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if u.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "account or password incorrect"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": &GetPersonalInfoResp{
			ID:       strconv.FormatInt(u.ID, 10),
			Account:  u.Account,
			UserName: u.Name,
		},
	})
}
