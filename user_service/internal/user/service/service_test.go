package service

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"user_service/internal/config"
	mock_user "user_service/internal/user/mock"
	"user_service/internal/user/model"
)

type usecase struct {
	Name      string
	UserInput model.User
	IsError   bool
}

const emptyData = ""

var (
	validUsers               []model.User
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
		Name:        "test_test",
		Login:       "test_test",
		Password:    "test_test",
		Description: "",
	}

	validUsers = []model.User{{
		ID:          "12345",
		Name:        "test",
		Login:       "test",
		Password:    "test",
		Description: "",
	}, {
		ID:          "678910",
		Name:        "test_test",
		Login:       "test_test",
		Password:    "test_test",
		Description: "",
	},
	}
	userEmpty = model.User{}

	code = m.Run()
}

func TestController_Create(t *testing.T) {
	var useCase = []usecase{
		{Name: "Failed to create user", UserInput: userEmpty, IsError: true},
		{Name: "Success to create user#1", UserInput: firstValidUserWithoutID, IsError: false},
		{Name: "Success to create user#2", UserInput: secondValidUserWithoutID, IsError: false},
	}

	mockCtrl := gomock.NewController(t)
	userRepoMock := mock_user.NewMockrepository(mockCtrl)
	cfg := config.Config{}

	unexpectedErr := errors.New("name is a vital field")

	userRepoMock.EXPECT().CreateUser(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, user model.User) error {
			if user.Name == emptyData {
				return unexpectedErr
			}
			return nil
		}).AnyTimes()

	srv := NewController(userRepoMock, &cfg)

	for _, us := range useCase {
		t.Run(us.Name, func(t *testing.T) {
			err := srv.Create(context.Background(), &us.UserInput)
			if us.IsError {
				assert.Error(t, err)

			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, us.UserInput.ID)
			}
		})
	}

}

func TestController_GetUser(t *testing.T) {
	var useCase = []usecase{
		{Name: "Failed to get user", UserInput: userEmpty, IsError: true},
		{Name: "Success to get user#1", UserInput: firstValidUser, IsError: false},
		{Name: "Success to get user#2", UserInput: secondValidUser, IsError: false},
	}

	mockCtrl := gomock.NewController(t)
	userRepoMock := mock_user.NewMockrepository(mockCtrl)
	cfg := config.Config{}

	userNotExist := errors.New("entity not found")

	userRepoMock.EXPECT().GetUser(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, id string) (model.User, error) {
			if id == emptyData {
				return model.User{}, userNotExist
			} else {
				return model.User{}, nil
			}
		}).AnyTimes()

	srv := NewController(userRepoMock, &cfg)

	for _, us := range useCase {
		t.Run(us.Name, func(t *testing.T) {
			_, err := srv.GetUser(context.Background(), us.UserInput.ID)
			if us.IsError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, us.UserInput.ID)
			}
		})
	}
}

func TestController_UpdateUser(t *testing.T) {
	var useCase = []usecase{
		{Name: "Failed to update user#1", UserInput: firstValidUserWithoutID, IsError: true},
		{Name: "Failed to update user#2", UserInput: secondValidUserWithoutID, IsError: true},
		{Name: "Failed to update user during DB transition", UserInput: secondValidUser, IsError: true},
		{Name: "Success to update user", UserInput: firstValidUser, IsError: false},
	}
	sqlUpdateErrId := secondValidUser.ID

	mockCtrl := gomock.NewController(t)
	userRepoMock := mock_user.NewMockrepository(mockCtrl)
	cfg := config.Config{}

	userNotExist := errors.New("entity not found")
	failedUpdate := errors.New("failed to update")

	userRepoMock.EXPECT().GetUser(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, id string) (model.User, error) {
			if id == emptyData {
				return model.User{}, userNotExist
			} else {
				return model.User{}, nil
			}
		}).AnyTimes()

	userRepoMock.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, user model.User) error {
			if user.ID == sqlUpdateErrId {
				return failedUpdate
			}
			return nil
		}).AnyTimes()

	srv := NewController(userRepoMock, &cfg)

	for _, us := range useCase {
		t.Run(us.Name, func(t *testing.T) {
			err := srv.UpdateUser(context.Background(), us.UserInput)
			if us.IsError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestController_DeleteUser(t *testing.T) {
	var useCase = []usecase{
		{Name: "Failed to delete user#01", UserInput: firstValidUserWithoutID, IsError: true},
		{Name: "Failed to delete user#02", UserInput: userEmpty, IsError: true},
		{Name: "Success to delete user", UserInput: firstValidUser, IsError: false},
	}

	mockCtrl := gomock.NewController(t)
	userRepoMock := mock_user.NewMockrepository(mockCtrl)
	cfg := config.Config{}

	userNotExist := errors.New("entity not found")

	userRepoMock.EXPECT().GetUser(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, id string) (model.User, error) {
			if id == emptyData {
				return model.User{}, userNotExist
			} else {
				return model.User{}, nil
			}
		}).AnyTimes()

	userRepoMock.EXPECT().DeleteUser(gomock.Any(), gomock.Any()).Times(1)

	srv := NewController(userRepoMock, &cfg)

	for _, us := range useCase {
		t.Run(us.Name, func(t *testing.T) {
			err := srv.DeleteUser(context.Background(), us.UserInput.ID)
			if us.IsError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestController_GetAllUsers(t *testing.T) {
	type usecase struct {
		Name    string
		Ctx     context.Context
		IsError bool
	}

	var useCase = []usecase{
		{Name: "Failed to get all users during DB transition", Ctx: nil, IsError: true},
		{Name: "Success to get all users", Ctx: context.Background(), IsError: false},
	}

	mockCtrl := gomock.NewController(t)
	userRepoMock := mock_user.NewMockrepository(mockCtrl)
	cfg := config.Config{}

	failedGetAllUsers := errors.New("failed to get all users")

	userRepoMock.EXPECT().GetAllUsers(gomock.Any()).DoAndReturn(
		func(ctx context.Context) ([]model.User, error) {
			if ctx == nil {
				return []model.User{}, failedGetAllUsers
			}
			return validUsers, nil
		}).Times(len(useCase))

	srv := NewController(userRepoMock, &cfg)

	for _, us := range useCase {
		t.Run(us.Name, func(t *testing.T) {
			users, err := srv.GetAllUsers(us.Ctx)
			if us.IsError {
				assert.Error(t, err)

			} else {
				assert.NoError(t, err)
				assert.Equal(t, users, validUsers)
			}
		})
	}
}

func TestController_Authorize(t *testing.T) {

	loginFirstValid := firstValidUser.Login
	loginSecondValid := secondValidUser.Login

	m := map[string]model.User{loginFirstValid: firstValidUser, loginSecondValid: secondValidUser}

	var useCase = []usecase{
		{Name: "Failed to authorize user", UserInput: userEmpty, IsError: true},
		{Name: "Success to authorize user#1", UserInput: m[loginFirstValid], IsError: false},
		{Name: "Success to authorize user#2", UserInput: m[loginSecondValid], IsError: false},
	}

	mockCtrl := gomock.NewController(t)
	userRepoMock := mock_user.NewMockrepository(mockCtrl)
	cfg := config.Config{}

	failedAuthorize := errors.New("failed to authorize user")

	userRepoMock.EXPECT().CheckAuth(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, login, password string) (model.User, error) {
			if login == emptyData && password == emptyData {
				return model.User{}, failedAuthorize
			} else {
				return m[login], nil
			}
		}).AnyTimes()

	srv := NewController(userRepoMock, &cfg)

	for _, us := range useCase {
		t.Run(us.Name, func(t *testing.T) {
			token, err := srv.Authorize(context.Background(), us.UserInput.Login, us.UserInput.Password)
			if us.IsError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, us.UserInput.Login, emptyData)
				assert.NotEqual(t, us.UserInput.Password, emptyData)
				var userClaim jwt.RegisteredClaims
				_, err = jwt.ParseWithClaims(token, &userClaim, func(token *jwt.Token) (interface{}, error) {
					return []byte("authorized"), nil
				})
				assert.Equal(t, us.UserInput.ID, userClaim.ID)
			}
		})
	}
}
