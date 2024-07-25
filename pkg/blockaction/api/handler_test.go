package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/reddtsai/goAPI/pkg/blockaction/storage"
	"github.com/reddtsai/goAPI/pkg/blockaction/storage/mock"
)

type TestBlockActionApi struct {
	suite.Suite
	ctrl        *gomock.Controller
	mockStorage *mock.MockIStorage

	TestApi *BlockActionApi
}

func TestBlockActionApiSuite(t *testing.T) {
	suite.Run(t, new(TestBlockActionApi))
}

func (t *TestBlockActionApi) SetupSuite() {
	t.ctrl = gomock.NewController(t.T())
	t.mockStorage = mock.NewMockIStorage(t.ctrl)
	api, err := NewBlockActionApi(SetStorage(t.mockStorage))
	if err != nil {
		t.FailNow(err.Error())
	}
	t.TestApi = api
}

func (t *TestBlockActionApi) Test_Signup_200() {
	payload := SignupReq{
		Account:  "testuser",
		Password: "abcd1234",
		UserName: "testuser",
	}
	body, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/signup", bytes.NewReader(body))

	t.mockStorage.EXPECT().IsExistUserAccount(payload.Account).Return(false, nil)
	t.mockStorage.EXPECT().CreateUser(gomock.Any()).Return(nil)

	t.TestApi.ServeHTTP(w, req)
	assert.Equal(t.T(), http.StatusOK, w.Code)
}

func (t *TestBlockActionApi) Test_Signup_400() {
	payload := SignupReq{
		Account:  "test",
		Password: "abcd1234",
		UserName: "testuser",
	}
	body, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/signup", bytes.NewReader(body))

	t.TestApi.ServeHTTP(w, req)
	assert.Equal(t.T(), http.StatusBadRequest, w.Code)
}

func (t *TestBlockActionApi) Test_Signup_409() {
	payload := SignupReq{
		Account:  "testuser",
		Password: "abcd1234",
		UserName: "testuser",
	}
	body, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/signup", bytes.NewReader(body))

	t.mockStorage.EXPECT().IsExistUserAccount(payload.Account).Return(true, nil)

	t.TestApi.ServeHTTP(w, req)
	assert.Equal(t.T(), http.StatusConflict, w.Code)
}

func (t *TestBlockActionApi) Test_Signin_200() {
	payload := SigninReq{
		Account:  "testuser",
		Password: "abcd1234",
	}
	body, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/signin", bytes.NewReader(body))

	s, err := signaturePwd(payload.Password)
	assert.Nil(t.T(), err)
	t.mockStorage.EXPECT().GetUserByAccount(payload.Account).Return(storage.UserTable{
		ID:      1,
		Account: "testuser",
		Secret:  s,
	}, nil)

	t.TestApi.ServeHTTP(w, req)
	assert.Equal(t.T(), http.StatusOK, w.Code)
}

func (t *TestBlockActionApi) Test_Signin_400() {
	payload := SigninReq{
		Account: "testuser",
	}
	body, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/signin", bytes.NewReader(body))

	t.TestApi.ServeHTTP(w, req)
	assert.Equal(t.T(), http.StatusBadRequest, w.Code)
}

func (t *TestBlockActionApi) Test_Signin_401() {
	payload := SigninReq{
		Account:  "testuser",
		Password: "abcd1234",
	}
	body, _ := json.Marshal(payload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/signin", bytes.NewReader(body))

	t.mockStorage.EXPECT().GetUserByAccount(payload.Account).Return(storage.UserTable{}, nil)

	t.TestApi.ServeHTTP(w, req)
	assert.Equal(t.T(), http.StatusUnauthorized, w.Code)
}

func (t *TestBlockActionApi) Test_GetPersonalInfo_200() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/user/personal-info", nil)
	token, err := genToken(&UserClaims{
		ID:      1,
		Account: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "blockaction",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(600 * time.Second)),
		},
	})
	assert.Nil(t.T(), err)
	bearer := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", bearer)

	t.mockStorage.EXPECT().GetUser(gomock.Any()).Return(storage.UserTable{
		ID:      1,
		Account: "testuser",
		Name:    "testuser",
	}, nil)

	t.TestApi.ServeHTTP(w, req)
	assert.Equal(t.T(), http.StatusOK, w.Code)
}

func (t *TestBlockActionApi) Test_GetPersonalInfo_401() {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/v1/user/personal-info", nil)

	t.TestApi.ServeHTTP(w, req)
	assert.Equal(t.T(), http.StatusUnauthorized, w.Code)
}
