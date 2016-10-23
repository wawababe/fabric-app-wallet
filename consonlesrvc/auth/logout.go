package authsrvc

import (
	"fmt"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"database/sql"
	"baas/app-wallet/consonlesrvc/database"
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

	var authRes *AuthResponse = new(AuthResponse)
	if !req.IsRequestValid(&res.AuthResponse) {
		authLogger.Warningf("request not valid: %#v", *req)
		res.UserUUID = ""
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
	fmt.Fprintf(w, "%s", string(resBytes))
}
