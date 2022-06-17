package userapi

import (
	"net/http"
	"time"

	"github.com/evertontomalok/distributed-system-go/internal/app"
)

func New(cfg app.Config) *Adapter {
	return &Adapter{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
		BaseUrl: cfg.UserApi.BaseUrl,
	}
}
