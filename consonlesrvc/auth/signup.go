package authsrvc

import (
	"fmt"
//	"database/sql"
	"github.com/julienschmidt/httprouter"
	_ "github.com/go-sql-driver/mysql"
	"baas/app-wallet/consonlesrvc/common"
	"net/http"
	"encoding/json"
	"database/sql"
	"baas/app-wallet/consonlesrvc/database"
)

type SignupRequest struct {
	username string
	password string
}

type SignupResponse struct {
	common.BaseResponse
	Session string `json:"sessionid"`
	AuthToken string `json:"authtoken"`
}

type Signup struct {
}


func (t *Signup) Post(req *SignupRequest)(*SignupResponse){
	var res *SignupResponse = new(SignupResponse)
	var err error
	var db *sql.DB = database.GetDB()
	var user *database.User = new(database.User)

	res.Status = "ok"

	if _, err = database.GetUserByName(db, req.username); err == nil {
		authLogger.Warningf("failed to sign up, user %s has existed", req.username)
		res.Status = "error"
		res.Message = "faile to signup duplicate user"
		return res
	}

	user.Username = req.username
	user.Password = req.password
	user.UserUUID = common.GenerateUUID()
	if _, err = database.AddUser(db, user); err != nil {
		authLogger.Errorf("failed to adduser %#v: %v", user, err)
		res.Status = "error"
		res.Message = "failed to signup, adduser error"
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




func SignupPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req *SignupRequest = new(SignupRequest)
	var res *SignupResponse
	var resBytes []byte
	var err error
	var t *Signup

	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()
	req.username = r.PostForm.Get("username")
	req.password = r.PostForm.Get("password")
	res = t.Post(req)

	resBytes, err = json.Marshal(*res)
	if err != nil {
		authLogger.Fatalf("failed to marshal response as []byte: %v", err)
	}
	fmt.Fprintf(w, "%s", string(resBytes))

}