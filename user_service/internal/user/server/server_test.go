package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
	"user_service/internal/config"
	"user_service/internal/user/mock"
	"user_service/internal/user/model"
)

type server struct {
	Name      string
	UserInput model.User
	IsError   bool
}

const emptyData = ""

var (
	firstValidUser           model.User
	secondValidUser          model.User
	firstValidUserWithoutID  model.User
	secondValidUserWithoutID model.User
	userEmpty                model.User
)

func TestMain(m *testing.M) {
	var code = 1
	defer func() { os.Exit(code) }()

	firstValidUserWithoutID = model.User{
		Name:        "test",
		Login:       "test",
		Password:    "test",
		Description: "",
	}

	secondValidUserWithoutID = model.User{
		Name:        "test_test",
		Login:       "test_test",
		Password:    "test_test",
		Description: "",
	}

	firstValidUser = model.User{
		ID:          "12345",
		Name:        "test",
		Login:       "test",
		Password:    "test",
		Description: "",
	}

	secondValidUser = model.User{
		ID:          "678910",
		Name:        "test-test",
		Login:       "test_test",
		Password:    "test_test",
		Description: "",
	}

	userEmpty = model.User{}

	code = m.Run()
}

func TestServer_Create(t *testing.T) {
	type tc struct {
		Name    string
		Input   string
		IsError bool
	}
	firstUser, err := json.Marshal(&firstValidUserWithoutID)
	if err != nil {
		return
	}

	secondUser, err := json.Marshal(&secondValidUserWithoutID)
	if err != nil {
		return
	}

	var tests = []tc{
		{Name: "Failed to create user", Input: "", IsError: true},
		{Name: "Failed to create user", Input: `{}`, IsError: true},
		{Name: "Success to create user#1", Input: string(firstUser), IsError: false},
		{Name: "Success to create user#2", Input: string(secondUser), IsError: false},
	}

	ctrl := gomock.NewController(t)
	mockController := mock.NewMockcontroller(ctrl)

	unexpectedErr := errors.New("name is a vital field")

	mockController.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, user *model.User) error {
			if user.Name == emptyData {
				return unexpectedErr
			}
			return nil
		}).AnyTimes()

	s := NewServer("", echo.New(), logrus.New(), mockController, &config.Config{})
	s.RegisterRoutes()

	srv := httptest.NewServer(s.r)
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			url := fmt.Sprintf("%s/user/create", srv.URL)
			req, err := http.NewRequest(echo.POST, url, strings.NewReader(test.Input))
			require.NoError(t, err)

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			client := http.Client{
				Timeout: time.Second,
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
			}

			e := echo.New()
			w := httptest.NewRecorder()
			ctx := e.NewContext(req, w)

			err = s.Create(ctx)
			if err != nil {
				assert.Equal(t, false, err == nil)
			}

			if test.IsError {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
			} else {
				assert.Equal(t, http.StatusCreated, resp.StatusCode)
			}
		})
	}

}

func TestServer_GetUser(t *testing.T) {

	IDFirstValid := firstValidUser.ID
	IDSecondValid := secondValidUser.ID

	m := map[string]model.User{firstValidUser.ID: firstValidUser, secondValidUser.ID: secondValidUser}

	tests := []server{
		{Name: "Failed to get user#1", UserInput: userEmpty, IsError: true},
		{Name: "Failed to get user#2", UserInput: firstValidUserWithoutID, IsError: true},
		{Name: "Success to get user", UserInput: m[IDFirstValid], IsError: false},
		{Name: "Success to get user", UserInput: m[IDSecondValid], IsError: false},
	}

	ctrl := gomock.NewController(t)
	mockController := mock.NewMockcontroller(ctrl)

	userNotExist := errors.New("entity not found")

	mockController.EXPECT().GetUser(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, id string) (model.User, error) {
			if id == emptyData {
				return model.User{}, userNotExist
			} else {
				return m[id], nil
			}
		}).AnyTimes()

	keyword := "test"

	s := NewServer("", echo.New(), logrus.New(), mockController, &config.Config{JWTKeyword: keyword})
	s.RegisterRoutes()

	srv := httptest.NewServer(s.r)
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			url := fmt.Sprintf("%s/user/%s", srv.URL, test.UserInput.ID)
			req, err := http.NewRequest(echo.GET, url, nil)
			require.NoError(t, err)

			now := time.Now()

			claims := jwt.RegisteredClaims{
				Issuer:    test.UserInput.ID,
				Subject:   "authorized",
				Audience:  nil,
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 1)),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now),
				ID:        test.UserInput.ID,
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString([]byte(keyword))
			require.NoError(t, err)

			req.Header.Set("Authorization", tokenString)

			client := http.Client{
				Timeout: time.Second,
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
			}

			e := echo.New()
			w := httptest.NewRecorder()
			ctx := e.NewContext(req, w)

			err = s.GetUser(ctx)
			if err != nil {
				assert.Equal(t, false, err == nil)
			}

			if test.IsError {
				assert.Equal(t, http.StatusNotFound, resp.StatusCode)
			} else {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			}
		})
	}
}

func TestServer_DeleteUser(t *testing.T) {
	IDFirstValid := firstValidUser.ID
	IDSecondValid := secondValidUser.ID

	m := map[string]model.User{firstValidUser.ID: firstValidUser, secondValidUser.ID: secondValidUser}

	tests := []server{
		{Name: "Failed to delete user#1", UserInput: userEmpty, IsError: true},
		{Name: "Failed to delete user#2", UserInput: firstValidUserWithoutID, IsError: true},
		{Name: "Success to delete user", UserInput: m[IDFirstValid], IsError: false},
		{Name: "Success to delete user", UserInput: m[IDSecondValid], IsError: false},
	}

	ctrl := gomock.NewController(t)
	mockController := mock.NewMockcontroller(ctrl)

	userNotExist := errors.New("entity not found")

	mockController.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, id string) error {
			if id == emptyData {
				return userNotExist
			} else {
				return nil
			}
		}).AnyTimes()

	keyword := "test"

	s := NewServer("", echo.New(), logrus.New(), mockController, &config.Config{JWTKeyword: keyword})
	s.RegisterRoutes()

	srv := httptest.NewServer(s.r)
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			url := fmt.Sprintf("%s/user/%s", srv.URL, test.UserInput.ID)
			req, err := http.NewRequest(echo.DELETE, url, nil)
			require.NoError(t, err)

			now := time.Now()

			claims := jwt.RegisteredClaims{
				Issuer:    test.UserInput.ID,
				Subject:   "authorized",
				Audience:  nil,
				ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 1)),
				NotBefore: jwt.NewNumericDate(now),
				IssuedAt:  jwt.NewNumericDate(now),
				ID:        test.UserInput.ID,
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString([]byte(keyword))
			require.NoError(t, err)

			req.Header.Set("Authorization", tokenString)

			client := http.Client{
				Timeout: time.Second,
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Error(err)
			}

			e := echo.New()
			w := httptest.NewRecorder()
			ctx := e.NewContext(req, w)

			err = s.DeleteUser(ctx)
			if err != nil {
				assert.Equal(t, false, err == nil)
			}

			if test.IsError {
				assert.Equal(t, http.StatusNotFound, resp.StatusCode)
			} else {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			}

		})
	}
}
