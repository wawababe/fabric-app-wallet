package transaction

import (
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"baas/app-wallet/consonlesrvc/auth"
	"baas/app-wallet/consonlesrvc/database"
	"database/sql"
	"fmt"
	"net/http"
)

type TransactionListRequest struct {
	authsrvc.AuthRequest
}

type TransactionListResponse struct {
	authsrvc.AuthResponse
	PayerTransactionList []*database.Transaction `json:"payertransactionlist"`
	PayeeTransactionList []*database.Transaction `json:"payeetransactionlist"`
}

type TransactionList struct {

}

func (c *TransactionList) post(req *TransactionListRequest)(*TransactionListResponse){
	var res *TransactionListResponse = new(TransactionListResponse)
	var err error
	var db *sql.DB = database.GetDB()

	if !req.IsRequestValid(&res.AuthResponse) {
		wtLogger.Warningf("request not valid: %#v", *req)
		res.UserUUID = ""
		return res
	}

	if res.PayerTransactionList, err = database.GetTransactionsByPayeruuid(db, res.UserUUID); err != nil {
		wtLogger.Errorf("failed to get transactions by payeruuid %s: %v", res.UserUUID, err)
		res.Status = "error"
		res.Message = fmt.Sprintf("failed to get transactions by payeruuid %s: %v", res.UserUUID, err)
		return res
	}

	if res.PayeeTransactionList, err = database.GetTransactionsByPayeeuuid(db, res.UserUUID); err != nil {
		wtLogger.Errorf("failed to get transactions by payeeuuid %s: %v", res.UserUUID, err)
		res.Status = "error"
		res.Message = fmt.Sprintf("failed to get transactions by payeeuuid %s: %v", res.UserUUID, err)
		return res
	}

	res.Status = "ok"
	res.UserUUID = res.UserUUID
	res.Message = "sucessed in listing transactions"
	return res
}

func TransactionListPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	var err error
	var req *TransactionListRequest = new(TransactionListRequest)
	var res *TransactionListResponse
	var resBytes []byte

	w.Header().Set("Content-Type", "application/json")

	if err = r.ParseForm(); err != nil {
		wtLogger.Fatalf("failed to parse request for url %s: %v", r.URL.Path, err)
	}

	req.Username = r.PostForm.Get("username")
	req.SessionID = r.PostForm.Get("sessionid")
	req.AuthToken = r.PostForm.Get("authtoken")
	wtLogger.Debugf("parsed request for url %s: %#v", r.URL.Path, req)

	var t TransactionList
	res = t.post(req)

	resBytes, err = json.Marshal(*res)
	if err != nil {
		wtLogger.Fatalf("failed to marshal response as []byte: %v", err)
	}
	fmt.Fprintf(w, "%s", string(resBytes))
}
