package account

import (
	"fmt"
	"net/http"
	"encoding/json"

	"baas/app-wallet/consonlesrvc/auth"
	"github.com/julienschmidt/httprouter"
	"database/sql"
	"baas/app-wallet/consonlesrvc/database"
	util "baas/app-wallet/consonlesrvc/common"
	"strings"
	"baas/app-wallet/consonlesrvc/wallet/task"
)

type CreateRequest struct {
	authsrvc.AuthRequest
	AccountName string `json:"accountname"`
}

type CreateResponse struct {
	authsrvc.AuthResponse
	AccountUUID string `json:"accountuuid,omitempty"`
	TaskUUID string `json:"taskuuid,omitempty"`
}

type Create struct {
}

func (c *Create) post(req *CreateRequest)(*CreateResponse){
	var res *CreateResponse = new(CreateResponse)
	var err error
	var db *sql.DB = database.GetDB()

	if !req.IsRequestValid(&res.AuthResponse) {
		wtLogger.Warningf("request not valid: %#v", *req)
		res.Status = "error"
		res.Message = util.ERROR_UNAUTHORIZED
		res.UserUUID = ""
		return res
	}
	if len(req.AccountName) == 0 {
		wtLogger.Warning("account name should not be empty")
		res.Status = "error"
		res.Message = util.ERROR_BADREQUEST +": account name should not be empty"
		res.UserUUID = ""
		return res
	}

	var account *database.Account = new(database.Account)
	var accountid = util.MD5string(res.UserUUID + req.AccountName)
	if _, err = database.GetAccountByAccountID(db, accountid); err == nil {
		wtLogger.Errorf("failed to create duplicate account %#v", account)
		res.Status = "error"
		res.Message = util.ERROR_BADREQUEST + ": failed to create duplicate account"
		res.UserUUID = ""
		return res
	}

	account.AccountUUID = util.GenerateUUID()
	account.UserUUID = res.UserUUID
	account.AccountName = req.AccountName
	account.AccountID = accountid
	account.Amount = 1000
	account.BC_TXUUID = ""
	account.Status = "pending"

	if _, err = database.AddAccount(db, account); err != nil {
		wtLogger.Errorf("failed adding account %#v: %v", account, err)
		res.Status = "error"
		res.Message = "failed adding account"
		return res
	}

	var crontask task.CronTask = new(task.AccountCreateTask)
	var taskuuid string
	taskuuid, err = crontask.Create(account.AccountUUID, task.TASK_TYPE_CREATE_ACCOUNT, task.TASK_STATE_INIT)
	if err != nil {
		wtLogger.Errorf("failed to create task for createaccount event: %v", err)
		res.Status = "error"
		res.Message = fmt.Sprintf("failed to create task for createaccount event: %v", err)
		res.AccountUUID = ""
		res.TaskUUID = ""
		account.Status = "failed"
		database.DeleteAccount(db, account)
		return res
	}


	res.Status = "ok"
	res.UserUUID = res.UserUUID
	res.Message = "sucessed in creating task for accountcreate event"
	res.AccountUUID = account.AccountUUID
	res.TaskUUID = taskuuid
	return res
}

func AccountCreatePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	var err error
	var req *CreateRequest = new(CreateRequest)
	var res *CreateResponse
	var resBytes []byte

	w.Header().Set("Content-Type", "application/json")

	if err = r.ParseForm(); err != nil {
		wtLogger.Fatalf("failed to parse request for url %s: %v", r.URL.Path, err)
	}

	req.Username = r.PostForm.Get("username")
	req.SessionID = r.PostForm.Get("sessionid")
	req.AuthToken = r.PostForm.Get("authtoken")
	req.AccountName = r.PostForm.Get("accountname")
	wtLogger.Debugf("parsed request for url %s: %#v", r.URL.Path, req)

	var t Create
	res = t.post(req)

	resBytes, err = json.Marshal(*res)
	if err != nil {
		wtLogger.Fatalf("failed to marshal response as []byte: %v", err)
	}
	if strings.Contains(res.Message, util.ERROR_UNAUTHORIZED){
		w.WriteHeader(http.StatusUnauthorized)
	}else if strings.Contains(res.Message, util.ERROR_BADREQUEST){
		w.WriteHeader(http.StatusBadRequest)
	}
	fmt.Fprintf(w, "%s", string(resBytes))
}