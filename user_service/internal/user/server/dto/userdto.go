package dto

type UserDto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Description string `json:"description"`
}
