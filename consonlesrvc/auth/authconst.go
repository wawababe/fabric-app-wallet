package authsrvc

import (
	"github.com/op/go-logging"
	"baas/app-wallet/consonlesrvc/common"
)

var authLogger *logging.Logger = common.NewLogger("authent")

const (
	SESSION_EXPIRATION_DAYS = 1
)