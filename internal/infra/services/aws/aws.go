package aws

import (
	"context"

	log "github.com/sirupsen/logrus"
)

func SendErrorToCloudWatch(c context.Context, Error error) {
	// TODO implement something mocked
	log.Infof("Error %+v sent to cloud watch. Context: %+v", Error, c)
}
