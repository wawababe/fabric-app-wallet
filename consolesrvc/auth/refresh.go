package authsrvc

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"baas/app-wallet/consolesrvc/database"
	"database/sql"
	util "baas/app-wallet/consolesrvc/common"
	"strings"
)

type RefreshRequest struct {
	AuthRequest
}

type RefreshResponse struct {
	AuthResponse
}


type Refresh struct {
}


func (t *Refresh) post(req *RefreshRequest) *RefreshResponse {
	var res *RefreshResponse = new(RefreshResponse)
	var err error
	var db *sql.DB = database.GetDB()

	if !req.IsRequestValid(&res.AuthResponse) {
		authLogger.Warningf("request not valid: %#v", *req)
		res.Status = "error"
		res.Message = util.ERROR_UNAUTHORIZED
		res.UserUUID = ""
		return res
	}

	var session *database.UserSession = new(database.UserSession)
	if session, err = database.GetUserSession(db, res.UserUUID, req.SessionID); err != nil {
		authLogger.Errorf("failed to refresh, can't getusersession by useruuid %s and sessionuuid %s", res.UserUUID, req.SessionID)
		res.Status = "error"
		res.Message = util.ERROR_UNAUTHORIZED + ": usersession not exist"
		return res
	}

	session.RefreshExpiredTimeByDays(SESSION_EXPIRATION_DAYS)
	database.UpdateUserSession(db, session)

	res.Status = "ok"
	res.UserUUID = res.UserUUID
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
		authLogger.Fatalf("failed to parse request for url %s: %v", r.URL.Path, err)
	}

	req.Username = r.PostForm.Get("username")
	req.SessionID = r.PostForm.Get("sessionid")
	req.AuthToken = r.PostForm.Get("authtoken")
	authLogger.Debugf("parsed request for url %s: %#v", r.URL.Path, req)

	var t Refresh
	res = t.post(req)

	resBytes, err = json.Marshal(*res)
	if err != nil {
		authLogger.Fatalf("failed to marshal response as []byte: %v", err)
	}

	if strings.Contains(res.Message, util.ERROR_UNAUTHORIZED){
		w.WriteHeader(http.StatusUnauthorized)
	}else if strings.Contains(res.Message, util.ERROR_BADREQUEST){
		w.WriteHeader(http.StatusBadRequest)
	}

	fmt.Fprintf(w, "%s", string(resBytes))
}
