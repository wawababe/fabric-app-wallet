package authsrvc

import (
	"baas/app-wallet/consonlesrvc/common"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"baas/app-wallet/consonlesrvc/database"
	"database/sql"
	"strings"
)

type LoginRequest struct {
	username string
	password string
}

type LoginResponse struct {
	common.BaseResponse
	Session   string `json:"sessionid,omitempty"`
	AuthToken string `json:"authtoken,omitempty"`
}

type Login struct {
}


func (t *Login) post(req *LoginRequest)(*LoginResponse){
	var err error
	var res *LoginResponse = new(LoginResponse)

	res.Status = "ok"

	var db *sql.DB = database.GetDB()
	var user *database.User = new(database.User)
	user, err = database.GetUserByName(db, req.username)
	if err != nil {
		authLogger.Errorf("failed to getuserbyname %s: %v", req.username, err)
		res.Status = "error"
		res.Message = common.ERROR_BADREQUEST + fmt.Sprintf(": user with name %s not exist", req.username)
		return res
	}

	if !strings.EqualFold(user.Password, req.password) {
		authLogger.Errorf("failed to login, wrong password %s for user %s", req.password, req.username)
		res.Status = "error"
		res.Message = common.ERROR_UNAUTHORIZED + ": wrong password"
		return res
	}

	var session *database.UserSession = new(database.UserSession)
	session.UserUUID = user.UserUUID
	session.SessionUUID = common.GenerateUUID()
	session.AddExpiredTimeByDays(SESSION_EXPIRATION_DAYS)
	//todo: make db as an interface, make use of reflect: db.Add(instance) ==> instance.Add()
	if _, err := database.AddUserSession(db, session); err != nil {
		authLogger.Errorf("failed to login, can't add user session %#v: %v", session, err)
		res.Status = "error"
		res.Message = "failed to login, can't generate new session"
		return res
	}
	res.Session = session.SessionUUID
	res.AuthToken = common.GenSessionToken(user.UserUUID, session.SessionUUID, user.Password)
	return res

}

// Router /auth/login: POST METHOD,
func LoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var err error
	var req *LoginRequest = new(LoginRequest)
	var res *LoginResponse
	var resBytes []byte

	w.Header().Set("Content-Type", "application/json")

	if err = r.ParseForm(); err != nil {
		authLogger.Fatalf("failed to parse request for url %s: %v", r.URL.Path, err)
	}

	req.username = r.PostForm.Get("username")
	req.password = r.PostForm.Get("password")
	//loginReq.username = r.FormValue("username")
	//loginReq.password = r.FormValue("password")
	authLogger.Debugf("parsed request for url %s: %#v", r.URL.Path, req)

	var t Login
	res = t.post(req)

	resBytes, err = json.Marshal(*res)
	if err != nil {
		authLogger.Fatalf("failed to marshal response as []byte: %v", err)
	}
	if strings.Contains(res.Message, common.ERROR_UNAUTHORIZED){
		w.WriteHeader(http.StatusUnauthorized)
	}else if strings.Contains(res.Message, common.ERROR_BADREQUEST){
		w.WriteHeader(http.StatusBadRequest)
	}
	fmt.Fprintf(w, "%s", string(resBytes))
}
