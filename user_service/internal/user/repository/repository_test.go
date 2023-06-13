package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"user_service/internal/config"
	"user_service/internal/user/model"
)

var userRepo *UserRepo

func TestMain(m *testing.M) {
	var code = 1
	defer func() { os.Exit(code) }()

	cfg := config.Config{DB: &config.DB{
		Driver:   "postgres",
		Password: "secretpassword",
		User:     "postgres",
		Name:     "postgres",
		SSLMode:  "disable",
	}}

	db, err := sqlx.Connect(cfg.DB.Driver, fmt.Sprintf("user=%s dbname=%s  sslmode=%s password= %s", cfg.DB.User,
		cfg.DB.Name, cfg.DB.SSLMode, cfg.DB.Password))

	if err != nil {
		log.Println(err)
		return
	}
	userRepo = NewUserRepo(db, cfg.DB)

	code = m.Run()
}

func TestUserRepo_CreateUser(t *testing.T) {
	userSuccess := model.User{
		ID:          uuid.New().String(),
		Name:        "test",
		Login:       "test",
		Password:    "test",
		Description: "",
	}

	var userNull model.User

	userWithoutID := model.User{
		Name:        "test",
		Login:       "test",
		Password:    "test",
		Description: "",
	}

	errSuccessUser := userRepo.CreateUser(context.Background(), userSuccess)
	errDublicateUser := userRepo.CreateUser(context.Background(), userSuccess)
	errNullUser := userRepo.CreateUser(context.Background(), userNull)
	errWithoutID := userRepo.CreateUser(context.Background(), userWithoutID)

	defer deleteUser(t, &userSuccess)

	assert.NoError(t, errSuccessUser)
	assert.Error(t, errDublicateUser)
	assert.Error(t, errNullUser)
	assert.Error(t, errWithoutID)

}

func TestUserRepo_GetUser(t *testing.T) {
	validId := "5ea844b6-8ea8-44c2-93b2-dffeeffd6255"
	unexpectedId := "5ea844b6-8e"

	userValid := model.User{
		ID:          validId,
		Name:        "test",
		Login:       "test",
		Password:    "test",
		Description: "",
	}

	createdUser(t, &userValid)
	defer deleteUser(t, &userValid)

	actualUser, actualErr := userRepo.GetUser(context.Background(), validId)
	_, unexpectedErr := userRepo.GetUser(context.Background(), unexpectedId)

	assert.NoError(t, actualErr)
	assert.Equal(t, actualUser.ID, validId)
	assert.Error(t, unexpectedErr)

}

func TestUserRepo_DeleteUser(t *testing.T) {
	userDel := model.User{
		ID:          uuid.New().String(),
		Name:        "test",
		Login:       "test",
		Password:    "test",
		Description: "",
	}

	unexpectedId := "515495"
	var userNull model.User

	createdUser(t, &userDel)

	err := userRepo.DeleteUser(context.Background(), userDel.ID)
	errUserNull := userRepo.DeleteUser(context.Background(), userNull.ID)
	errUnexpected := userRepo.DeleteUser(context.Background(), unexpectedId)

	assert.NoError(t, err)
	assert.Error(t, errUserNull)
	assert.Error(t, errUnexpected)

}

func TestUserRepo_UpdateUser(t *testing.T) {
	user := model.User{
		ID:          uuid.New().String(),
		Name:        "test",
		Login:       "test",
		Password:    "test",
		Description: "",
	}

	userToUpdate := model.User{
		ID:          user.ID,
		Name:        "testUpdate",
		Login:       "testUpdate",
		Password:    "testUpdate",
		Description: "",
	}

	toUpdateUserNull := model.User{}

	createdUser(t, &user)

	defer deleteUser(t, &userToUpdate)

	errValidUser := userRepo.UpdateUser(context.Background(), userToUpdate)
	errUserNull := userRepo.UpdateUser(context.Background(), toUpdateUserNull)

	updatedUser, err := getUser(t, &userToUpdate)

	assert.NoError(t, err)
	assert.NoError(t, errValidUser)
	assert.Error(t, errUserNull)
	assert.Equal(t, updatedUser.Name, userToUpdate.Name)
	assert.Equal(t, updatedUser.Password, userToUpdate.Password)
	assert.Equal(t, updatedUser.Login, userToUpdate.Login)

}

func TestUserRepo_CheckAuth(t *testing.T) {
	validUser := model.User{
		ID:          uuid.New().String(),
		Name:        "test",
		Login:       "test",
		Password:    "test",
		Description: "",
	}

	createdUser(t, &validUser)
	defer deleteUser(t, &validUser)

	auth, err := userRepo.CheckAuth(context.Background(), validUser.Login, validUser.Password)
	_, errNonexistent := userRepo.CheckAuth(context.Background(), "Nonexistent login", "Nonexistent password")

	assert.NoError(t, err)
	assert.Equal(t, auth.Login, validUser.Login)
	assert.Equal(t, auth.Password, validUser.Password)
	assert.Error(t, errNonexistent)
}

func TestUserRepo_GetAllUsers(t *testing.T) {
	getUser1 := model.User{
		ID:          uuid.New().String(),
		Name:        "test",
		Login:       "test",
		Password:    "test",
		Description: "",
	}

	createdUser(t, &getUser1)
	defer deleteUser(t, &getUser1)

	users, err := userRepo.GetAllUsers(context.Background())

	assert.NoError(t, err)
	assert.NotEmpty(t, users)
}

func createdUser(t *testing.T, user *model.User) {
	err := userRepo.CreateUser(context.Background(), *user)
	assert.NoError(t, err)

}

func deleteUser(t *testing.T, tData *model.User) {
	err := userRepo.DeleteUser(context.Background(), tData.ID)
	assert.NoError(t, err)
}

func getUser(t *testing.T, tData *model.User) (*model.User, error) {
	user, err := userRepo.GetUser(context.Background(), tData.ID)
	assert.NoError(t, err)
	return &user, nil
}
