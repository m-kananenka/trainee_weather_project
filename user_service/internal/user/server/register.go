package server

import (
	"user_service/internal/middleware"
)

func (s *Server) RegisterRoutes() {
	s.r.Use(middleware.Logger(s.logger))

	s.r.POST("/user/create", s.Create)
	s.r.GET("/user/auth", s.Authorize)

	e := s.r.Group("/user")

	e.Use(middleware.Auth(s.cfg.JWTKeyword, s.logger))

	e.GET("/:userID", s.GetUser)
	e.PUT("/:userID", s.UpdateUser)
	e.DELETE("/:userID", s.DeleteUser)
	e.GET("", s.GetAllUsers)

}
