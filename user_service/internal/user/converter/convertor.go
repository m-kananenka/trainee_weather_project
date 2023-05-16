package converter

import (
	"github.com/jinzhu/copier"
	"user_service/internal/user/model"
	"user_service/internal/user/server/dto"
)

func UserToUserDTO(from *model.User) *dto.UserDto {
	target := &dto.UserDto{}

	err := copier.Copy(target, from)
	if err != nil {
		return nil
	}
	return target
}
func UserDtoToUser(from *dto.UserDto) *model.User {
	target := &model.User{}

	err := copier.Copy(target, from)
	if err != nil {
		return nil
	}
	return target
}
