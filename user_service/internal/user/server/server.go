package server

import (
	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"user_service/api/pb"
	"user_service/internal/config"
	"user_service/internal/errorstore"
	"user_service/internal/user/converter"
	"user_service/internal/user/model"
	"user_service/internal/user/server/dto"
)

type controller interface {
	Create(context context.Context, user *model.User) (*model.User, error)
	GetUser(ctx context.Context, id string) (model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	DeleteUser(ctx context.Context, id string) error
	GetAllUsers(ctx context.Context) ([]model.User, error)
	Authorize(ctx context.Context, login, password string) (string, error)
	Get(ctx context.Context, req *pb.Request) (*pb.Response, error)
}

type Server struct {
	listenURI string
	logger    *logrus.Logger
	r         *echo.Echo
	c         controller
	cfg       *config.Config
	client    *grpc.Server
}

func NewServer(listenURI string, r *echo.Echo, logger *logrus.Logger, c controller, cfg *config.Config) *Server {
	return &Server{
		listenURI: listenURI,
		logger:    logger,
		r:         r,
		c:         c,
		cfg:       cfg,
		client:    grpc.NewServer(),
	}
}

func (s *Server) Register() {
	s.client.RegisterService(&pb.UserService_ServiceDesc, s.c)
}
func (s *Server) StartGRPC() {
	l, err := net.Listen("tcp", "localhost:8085")
	if err != nil {
		s.logger.Fatal(err)
	}

	if err = s.client.Serve(l); err != nil {
		s.logger.Fatal(err)
	}
}

func (s *Server) StartRouter() {
	srv := http.Server{
		Addr:    s.listenURI,
		Handler: s.r,
	}
	s.logger.Info("server is running....")
	err := srv.ListenAndServe()
	if err != nil {
		s.logger.Fatal(err)

	}
}

func (s *Server) Create(ctx echo.Context) error {
	var userDto dto.UserDto
	err := ctx.Bind(&userDto)
	if err != nil {
		s.logger.Info("could not decode data:", err)
		return echo.NewHTTPError(http.StatusBadRequest, errorstore.BadRequest(err))
	}

	createdUser, err := s.c.Create(ctx.Request().Context(), converter.UserDtoToUser(&userDto))
	if err != nil {
		s.logger.Info("error parsing context during creating:", err)
		return echo.NewHTTPError(http.StatusBadRequest, errorstore.BadRequest(err))
	} else {
		return echo.NewHTTPError(http.StatusCreated, createdUser)
	}

}

func (s *Server) GetUser(ctx echo.Context) error {
	userID := ctx.Param("userID")

	user, err := s.c.GetUser(ctx.Request().Context(), userID)
	if err != nil {
		s.logger.Info("could not get user:", err)
		return echo.NewHTTPError(http.StatusBadRequest, errorstore.EntityNotFound(err))
	} else {
		return echo.NewHTTPError(http.StatusOK, converter.UserToUserDTO(&user))
	}
}

func (s *Server) UpdateUser(ctx echo.Context) error {
	var userDto dto.UserDto
	err := ctx.Bind(&userDto)
	if err != nil {
		s.logger.Info("could not decode data:", err)
		return echo.NewHTTPError(http.StatusBadRequest, errorstore.BadRequest(err))
	}
	userDto.ID = ctx.Param("userID")

	err = s.c.UpdateUser(ctx.Request().Context(), *converter.UserDtoToUser(&userDto))
	if err != nil {
		s.logger.Info("error parsing context during updating:", err)
		return echo.NewHTTPError(http.StatusBadRequest, errorstore.EntityNotFound(err))
	} else {
		return echo.NewHTTPError(http.StatusOK, userDto)
	}
}

func (s *Server) DeleteUser(ctx echo.Context) error {
	userID := ctx.Param("userID")

	err := s.c.DeleteUser(ctx.Request().Context(), userID)
	if err != nil {
		s.logger.Info("could not delete user:", err)
		return echo.NewHTTPError(http.StatusBadRequest, errorstore.EntityNotFound(err))
	} else {
		return echo.NewHTTPError(http.StatusOK)
	}

}

func (s *Server) GetAllUsers(ctx echo.Context) error {
	var userDto []*dto.UserDto
	users, err := s.c.GetAllUsers(ctx.Request().Context())

	if err != nil {
		s.logger.Info("could not get all users:", err)
		return echo.NewHTTPError(http.StatusBadRequest, errorstore.EntityNotFound(err))
	} else {
		err := copier.Copy(&userDto, users)
		if err != nil {
			return err
		}
		return echo.NewHTTPError(http.StatusOK, userDto)
	}
}

func (s *Server) Authorize(ctx echo.Context) error {
	var userDto dto.UserDto
	err := ctx.Bind(&userDto)
	if err != nil {
		s.logger.Info("could not decode data:", err)
		return echo.NewHTTPError(http.StatusBadRequest, errorstore.BadRequest(err))
	}

	token, err := s.c.Authorize(ctx.Request().Context(), userDto.Login, userDto.Password)
	if err != nil {
		s.logger.Info("error parsing context during authorization:", err)
		return echo.NewHTTPError(http.StatusBadRequest, errorstore.BadRequest(err))
	} else {
		return echo.NewHTTPError(http.StatusOK, token)
	}

}
