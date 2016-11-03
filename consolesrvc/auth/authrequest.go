package authsrvc

import (
	"baas/app-wallet/consonlesrvc/database"
	"strings"
	"database/sql"
	"baas/app-wallet/consonlesrvc/common"
)

type AuthRequest struct {
	Username string `json:"username"`
	SessionID string `json:"sessionid"`
	AuthToken string `json:"authtoken"`
}

type AuthResponse struct {
	common.BaseResponse
	UserUUID string `json:"useruuid,omitempty"`
}

//IsAuthRequestValid: check whether request req is valid; set the response res
func (req *AuthRequest) IsRequestValid(res *AuthResponse)(bool){
	var err error
	var db *sql.DB = database.GetDB()
	var user *database.User = new(database.User)

	if user, err = database.GetUserByName(db, req.Username); err != nil {
		authLogger.Warningf("failed to validate request, user %s not exists", req.Username)
		res.Status = "error"
		res.Message = "failed to validate request, user not exist"
		return false
	}


	var session *database.UserSession = new(database.UserSession)

	if !strings.EqualFold(req.AuthToken, common.GenSessionToken(user.UserUUID, req.SessionID, user.Password)){
		authLogger.Infof("failed to validate request, wrong authtoken: %s", req.AuthToken)
		res.Status = "error"
		res.Message = "failed to validate request, wrong authtoken"
		return false
	}


	if session, err = database.GetUserSession(db, user.UserUUID, req.SessionID); err != nil {
		authLogger.Errorf("failed to validate request, can't getusersession by useruuid %s and sessionuuid %s", user.UserUUID, req.SessionID)
		res.Status = "error"
		res.Message = "failed to validate request, can't getusersession"
		return false
	}

	if session.IsExpired() {
		authLogger.Errorf("failed to validate request, session has expired in %s", session.ExpiredAt)
		res.Status = "error"
		res.Message = "failed to validate request, session has expired at " + session.ExpiredAt
		return false
	}
	res.Status = "ok"
	res.UserUUID = user.UserUUID
	return true
}