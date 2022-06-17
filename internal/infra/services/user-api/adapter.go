package userapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/evertontomalok/distributed-system-go/internal/domain/core/dto"
	log "github.com/sirupsen/logrus"
)

var UserAdapter *Adapter

type Adapter struct {
	Client  *http.Client
	BaseUrl string
}

var endpoints map[string]string = map[string]string{
	"user": "user/%s",
}

func (a *Adapter) GetUserStatus(userId string) dto.UserResponse {
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
		log.Panic(err)
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	userResponse := dto.UserResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userResponse)

	return userResponse
}
