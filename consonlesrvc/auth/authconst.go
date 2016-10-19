package authsrvc

import "github.com/op/go-logging"

var authLogger *logging.Logger = logging.MustGetLogger("authent")

const (
	SESSION_EXPIRATION_DAYS = 1
)