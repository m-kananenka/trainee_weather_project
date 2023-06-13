package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strings"
	pb2 "telegram_service/cmd/user/pb"
)

type AuthService struct {
	m map[int64]bool
}

func NewAuthService() *AuthService {
	return &AuthService{
		m: map[int64]bool{},
	}
}

func (a *AuthService) CheckAuth(id int64) bool {
	v, ok := a.m[id]
	if ok {
		return v
	} else {
		return false
	}
}

func (a *AuthService) Auth(text string, id int64) bool {
	login, password, err := SplitString(text)
	if err != nil {
		a.m[id] = false
		return false
	}
	conn, err := grpc.Dial("localhost:8085", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect: %v", err)
		a.m[id] = false
		return false
	}
	defer conn.Close()

	userClient := pb2.NewUserServiceClient(conn)

	req := &pb2.Request{
		Login:    login,
		Password: password,
	}

	res, err := userClient.Get(context.Background(), req)
	if err != nil {
		log.Printf("Failed to call authorization: %v", err)
		return false
	}
	a.m[id] = res.GetResponse()

	return res.GetResponse()
}

func SplitString(s string) (login string, password string, err error) {
	split := strings.Split(s, " ")

	if len(split) != 2 {
		return "", "", fmt.Errorf("failed to split string %w", err)
	}
	return split[0], split[1], nil
}
