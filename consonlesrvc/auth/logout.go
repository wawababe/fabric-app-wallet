package authsrvc

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"baas/app-wallet/consonlesrvc/common"
	"encoding/json"
	"database/sql"
	"baas/app-wallet/consonlesrvc/database"
)

type LogoutRequest struct{
	AuthRequest
}

func (t *LogoutRequest) CopyTo(dst *AuthRequest){
	dst.SessionID = t.SessionID
	dst.Username = t.Username
	dst.AuthToken = t.AuthToken
}

type LogoutResponse struct{
	common.BaseResponse
}

func (t *LogoutResponse) CopyFrom(src *AuthResponse){
	t.Status = src.Status
	t.Message = src.Message
}

type Logout struct {
}


func (t *Logout) Post(req *LogoutRequest) *LogoutResponse {
	var res *LogoutResponse = new(LogoutResponse)
	var err error
	var db *sql.DB = database.GetDB()

	var authReq *AuthRequest = new(AuthRequest)
	var authRes *AuthResponse = new(AuthResponse)
	req.CopyTo(authReq)
	if !authReq.IsAuthRequestValid(authRes) {
		authLogger.Warningf("request not valid: %#v", *authReq)
		res.CopyFrom(authRes)
		return res
	}

	var session *database.UserSession = new(database.UserSession)
	if session, err = database.GetUserSession(db, authRes.UserUUID, req.SessionID); err != nil {
		authLogger.Errorf("failed to logout, can't getusersession by useruuid %s and sessionuuid %s", authRes.UserUUID, req.SessionID)
		res.Status = "error"
		res.Message = "failed to refresh, can't getusersession"
		return res
	}

	//todo: need to add deleted field for usersession table; then update this field as 1(deleted)
	database.UpdateUserSession(db, session)

	res.Status = "ok"
	res.Message = "sucessed in logging out"
	return res
}


func LogoutPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	var req *LogoutRequest = new(LogoutRequest)
	var res *LogoutResponse
	var resBytes []byte
	var err error
	var t *Logout

	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()
	req.Username = r.PostForm.Get("username")
	req.SessionID = r.PostForm.Get("sessionid")
	req.AuthToken = r.PostForm.Get("authtoken")
	res = t.Post(req)

	resBytes, err = json.Marshal(*res)
	if err != nil {
		authLogger.Fatalf("failed to marshal response as []byte: %v", err)
	}
	fmt.Fprintf(w, "%s", string(resBytes))
}
