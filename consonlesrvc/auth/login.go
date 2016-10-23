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
	Session   string `json:"sessionid"`
	AuthToken string `json:"authtoken"`
}

type Login struct {
}


func (t *Login) post(loginReq *LoginRequest)(*LoginResponse){
	var err error
	var loginRes *LoginResponse = new(LoginResponse)

	loginRes.Status = "ok"

	var db *sql.DB = database.GetDB()
	var user *database.User = new(database.User)
	user, err = database.GetUserByName(db, loginReq.username)
	if err != nil {
		authLogger.Errorf("failed to getuserbyname %s: %v", loginReq.username, err)
		loginRes.Status = "error"
		loginRes.Message = "failed to get user: " + err.Error()
		return loginRes
	}

	if !strings.EqualFold(user.Password, loginReq.password) {
		authLogger.Errorf("failed to login, wrong password %s for user %s", loginReq.password, loginReq.username)
		loginRes.Status = "error"
		loginRes.Message = "failed to login, wrong password"
		return loginRes
	}

	var session *database.UserSession = new(database.UserSession)
	session.UserUUID = user.UserUUID
	session.SessionUUID = common.GenerateUUID()
	session.AddExpiredTimeByDays(SESSION_EXPIRATION_DAYS)
	//todo: make db as an interface, make use of reflect: db.Add(instance) ==> instance.Add()
	if _, err := database.AddUserSession(db, session); err != nil {
		authLogger.Errorf("failed to login, can't add user session %#v: %v", session, err)
		loginRes.Status = "error"
		loginRes.Message = "failed to login, can't generate new session"
		return loginRes
	}
	loginRes.Session = session.SessionUUID
	loginRes.AuthToken = common.GenSessionToken(user.UserUUID, session.SessionUUID, user.Password)
	return loginRes

}

// Router /auth/login: POST METHOD,
func LoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var err error
	var loginReq *LoginRequest = new(LoginRequest)
	var loginRes *LoginResponse
	var resBytes []byte

	w.Header().Set("Content-Type", "application/json")

	if err = r.ParseForm(); err != nil {
		authLogger.Fatalf("failed to parse request for url %s: %v", r.URL.Path, err)
	}

	loginReq.username = r.PostForm.Get("username")
	loginReq.password = r.PostForm.Get("password")
	//loginReq.username = r.FormValue("username")
	//loginReq.password = r.FormValue("password")
	authLogger.Debugf("parsed request for url %s: %#v", r.URL.Path, loginReq)

	var t Login
	loginRes = t.post(loginReq)

	resBytes, err = json.Marshal(*loginRes)
	if err != nil {
		authLogger.Fatalf("failed to marshal response as []byte: %v", err)
	}
	fmt.Fprintf(w, "%s", string(resBytes))
}
