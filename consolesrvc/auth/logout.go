package authsrvc

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"database/sql"
	"baas/app-wallet/consonlesrvc/database"
	util "baas/app-wallet/consonlesrvc/common"
	"strings"
)

type LogoutRequest struct{
	AuthRequest
}

type LogoutResponse struct{
	AuthResponse
}

type Logout struct {
}


func (t *Logout) post(req *LogoutRequest) *LogoutResponse {
	var res *LogoutResponse = new(LogoutResponse)
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
		authLogger.Errorf("failed to logout, can't getusersession by useruuid %s and sessionuuid %s", res.UserUUID, req.SessionID)
		res.Status = "error"
		res.Message = "failed to logout, can't getusersession"
		return res
	}

	//todo: need to add deleted field for usersession table; then update this field as 1(deleted)
	if _, err = database.DeleteUserSession(db, session); err != nil {
		authLogger.Errorf("failed to logout, can't delete usersession %#v: %v", *session, err)
		res.Status = "error"
		res.Message = "failed to logout: " + err.Error()
	}

	res.Status = "ok"
	res.Message = "sucessed in logging out"
	return res
}

// route: /auth/logout, method: POST
func LogoutPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params){
	var req *LogoutRequest = new(LogoutRequest)
	var res *LogoutResponse
	var resBytes []byte
	var err error


	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()
	req.Username = r.PostForm.Get("username")
	req.SessionID = r.PostForm.Get("sessionid")
	req.AuthToken = r.PostForm.Get("authtoken")

	var t Logout
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
