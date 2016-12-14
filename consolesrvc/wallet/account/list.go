package account

import (
	"baas/app-wallet/consolesrvc/auth"
	"baas/app-wallet/consolesrvc/database"
	"github.com/julienschmidt/httprouter"
	util "baas/app-wallet/consolesrvc/common"
	"encoding/json"
	"database/sql"
	"net/http"
	"fmt"
	"strings"
)

type AccountListRequest struct {
	authsrvc.AuthRequest
}

type AccountListResponse struct {
	authsrvc.AuthResponse
	AccountList []*database.Account `json:"accountlist,omitempty"`
}

type AccountList struct {

}

func (c *AccountList) post(req *AccountListRequest)(*AccountListResponse){
	var res *AccountListResponse = new(AccountListResponse)
	var err error
	var db *sql.DB = database.GetDB()

	if !req.IsRequestValid(&res.AuthResponse) {
		wtLogger.Warningf("request not valid: %#v", *req)
		res.Status = "error"
		res.Message = util.ERROR_UNAUTHORIZED
		res.UserUUID = ""
		return res
	}

	if res.AccountList, err = database.GetAccountsByUseruuid(db, res.UserUUID); err != nil {
		wtLogger.Errorf("failed to get accounts by useruuid %s: %v", res.UserUUID, err)
		res.Status = "error"
		res.Message = fmt.Sprintf("failed to get accounts by useruuid %s: %v", res.UserUUID, err)
		return res
	}

	res.Status = "ok"
	res.UserUUID = res.UserUUID
	res.Message = "sucessed in listing accounts"
	return res
}

func AccountListPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	var err error
	var req *AccountListRequest = new(AccountListRequest)
	var res *AccountListResponse
	var resBytes []byte

	w.Header().Set("Content-Type", "application/json")

	if err = r.ParseForm(); err != nil {
		wtLogger.Fatalf("failed to parse request for url %s: %v", r.URL.Path, err)
	}

	req.Username = r.PostForm.Get("username")
	req.SessionID = r.PostForm.Get("sessionid")
	req.AuthToken = r.PostForm.Get("authtoken")
	wtLogger.Debugf("parsed request for url %s: %#v", r.URL.Path, req)

	var t AccountList
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
