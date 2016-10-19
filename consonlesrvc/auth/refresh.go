package authsrvc

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"baas/app-wallet/consonlesrvc/database"
	"database/sql"
)

type RefreshRequest struct {
	AuthRequest
}

func (t *RefreshRequest) CopyTo(dst *AuthRequest){
	dst.SessionID = t.SessionID
	dst.Username = t.Username
	dst.AuthToken = t.AuthToken
}

type RefreshResponse struct {
	AuthResponse
}

func (t *RefreshResponse) CopyFrom(src *AuthResponse){
	t.Status = src.Status
	t.Message = src.Message
	t.UserUUID = src.UserUUID
}

type Refresh struct {
}


func (t *Refresh) Post(req *RefreshRequest) *RefreshResponse {
	var res *RefreshResponse = new(RefreshResponse)
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
		authLogger.Errorf("failed to refresh, can't getusersession by useruuid %s and sessionuuid %s", authRes.UserUUID, req.SessionID)
		res.Status = "error"
		res.Message = "failed to refresh, can't getusersession"
		return res
	}

	session.RefreshExpiredTimeByDays(SESSION_EXPIRATION_DAYS)
	database.UpdateUserSession(db, session)

	res.Status = "ok"
	res.UserUUID = authRes.UserUUID
	res.Message = "sucessed in refreshing session"
	return res
}

func RefreshPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	var err error
	var req *RefreshRequest = new(RefreshRequest)
	var res *RefreshResponse
	var resBytes []byte

	w.Header().Set("Content-Type", "application/json")

	if err = r.ParseForm(); err != nil {
		authLogger.Fatalf("failed to parse request, router /auth/login: %v", err)
	}

	req.Username = r.PostForm.Get("username")
	req.SessionID = r.PostForm.Get("sessionid")
	req.AuthToken = r.PostForm.Get("authtoken")
	authLogger.Debugf("parsed request for /auth/refresh: %#v", req)

	var t Refresh
	res = t.Post(req)

	resBytes, err = json.Marshal(*res)
	if err != nil {
		authLogger.Fatalf("failed to marshal response as []byte: %v", err)
	}
	fmt.Fprintf(w, "%s", string(resBytes))
}
