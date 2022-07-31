package dto

type UserResponse struct {
	UserId  string `json:"user_id"`
	IsValid bool   `json:"is_valid"`
}

func (u UserResponse) ReturnWithError(err error) (UserResponse, error) {
	return UserResponse{}, err
}
