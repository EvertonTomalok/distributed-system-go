package userapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/evertontomalok/distributed-system-go/internal/core/dto"
)

var UserAdapter *Adapter

type Adapter struct {
	Client  *http.Client
	BaseUrl string
}

var endpoints map[string]string = map[string]string{
	"user": "user/%s",
}

func (a *Adapter) GetUserStatus(userId string) (dto.UserResponse, error) {
	userResponse := dto.UserResponse{}

	endpoint := fmt.Sprintf(endpoints["user"], userId)
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"%s/%s",
			strings.TrimRight(a.BaseUrl, "/"),
			endpoint,
		),
		nil,
	)

	if err != nil {
		return userResponse.ReturnWithError(err)
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return userResponse.ReturnWithError(err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&userResponse)

	if err != nil {
		return userResponse.ReturnWithError(err)
	}

	return userResponse, nil
}
