package account

import (
	"baas/app-wallet/consonlesrvc/auth"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"net/http"
	"fmt"
	"baas/app-wallet/consonlesrvc/database"
	"database/sql"
	"strconv"
	"strings"
	util "baas/app-wallet/consonlesrvc/common"
	"baas/app-wallet/consonlesrvc/wallet/task"
)

type TransferRequest struct {
	authsrvc.AuthRequest
	PayerAccountID string `json:"payeraccountid"`
	PayeeAccountID string `json:"payeeaccountid"`
	Amount int64
}

type TransferResponse struct {
	authsrvc.AuthResponse
	TaskUUID string `json:"taskuuid,omitempty"`
}


type Transfer struct {

}

func (t *Transfer) post(req *TransferRequest) (*TransferResponse) {
	var res *TransferResponse = new(TransferResponse)
	var err error
	var db *sql.DB = database.GetDB()

	if !req.IsRequestValid(&res.AuthResponse) {
		wtLogger.Warningf("request not valid: %#v", *req)
		res.Status = "error"
		res.Message = util.ERROR_UNAUTHORIZED
		res.UserUUID = ""
		return res
	}

	if req.Amount <= 0 {
		wtLogger.Errorf("illegal to pay %d in transfer event", req.Amount)
		res.Status = "error"
		res.Message = util.ERROR_BADREQUEST + fmt.Sprintf(": illegal to pay amount %d", req.Amount)
		res.UserUUID = ""
		return res
	}

	var payerAccount *database.Account = new(database.Account)
	var payeeAccount *database.Account = new(database.Account)
	if payerAccount, err = database.GetAccountByAccountID(db, req.PayerAccountID); err != nil {
		wtLogger.Errorf("failed to get payer account by accountid", req.PayerAccountID)
		res.Status = "error"
		res.Message = util.ERROR_BADREQUEST + fmt.Sprintf(": failed to get payer account by accountid %s", req.PayerAccountID)
		return res
	}
	if !strings.EqualFold(res.UserUUID, payerAccount.UserUUID) {
		wtLogger.Errorf("failed to validate, user %s with useruuid %s not have account %#v", req.Username, res.UserUUID, payerAccount)
		res.Status = "error"
		res.Message = util.ERROR_UNAUTHORIZED + fmt.Sprintf(": request not valid, user %s not hava account %s", req.Username, req.PayerAccountID)
		return res
	}

	if payeeAccount, err = database.GetAccountByAccountID(db, req.PayeeAccountID); err != nil {
		wtLogger.Errorf("failed to get payee account %s: %v", req.PayeeAccountID, err)
		res.Status = "error"
		res.Message = util.ERROR_BADREQUEST + fmt.Sprintf(": failed to get payee account %s", req.PayeeAccountID)
		return res
	}

	if payerAccount.Amount < req.Amount {
		wtLogger.Errorf("failed to pay, payer's account only has %d, not enough to pay %d", payerAccount.Amount, req.Amount)
		res.Status = "error"
		res.Message = util.ERROR_BADREQUEST + ": account residual amount is not enough to pay"
		return res
	}


	var tx *database.Transaction = new(database.Transaction)
	tx.TxUUID = util.GenerateUUID()
	tx.PayerUUID = payerAccount.UserUUID
	tx.PayeeUUID = payeeAccount.UserUUID
	tx.PayerAccountID = payerAccount.AccountID
	tx.PayeeAccountID = payeeAccount.AccountID
	tx.Amount = req.Amount
	tx.Status = "pending"


	if _, err = database.AddTransaction(db, tx); err != nil {
		wtLogger.Errorf("failed adding transaction %#v: %v", tx, err)
		res.Status = "error"
		res.Message = "failed adding transaction"
		return res
	}

	var crontask task.CronTask = new(task.AccountTransferTask)
	var taskuuid string
	taskuuid, err = crontask.Create(tx.TxUUID, task.TASK_TYPE_TRANSFER, task.TASK_STATE_INIT)
	if err != nil {
		wtLogger.Errorf("failed to create task for createaccount event: %v", err)
		res.Status = "error"
		res.Message = fmt.Sprintf("failed to create task for createaccount event: %v", err)
		res.TaskUUID = ""
		tx.Status = "failed"
		database.DeleteTransaction(db, tx)
		return res
	}

	res.Status = "ok"
	res.Message = "sucessed in creating account"
	res.UserUUID = res.UserUUID
	res.TaskUUID = taskuuid
	return res
}


func TransferPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	var err error
	var req *TransferRequest = new(TransferRequest)
	var res *TransferResponse
	var resBytes []byte

	w.Header().Set("Content-Type", "application/json")

	if err = r.ParseForm(); err != nil {
		wtLogger.Fatalf("failed to parse request for url %s: %v", r.URL.Path, err)
	}

	req.Username = r.PostForm.Get("username")
	req.SessionID = r.PostForm.Get("sessionid")
	req.AuthToken = r.PostForm.Get("authtoken")
	req.PayerAccountID = r.PostForm.Get("payeraccountid")
	req.PayeeAccountID = r.PostForm.Get("payeeaccountid")
	req.Amount, _ = strconv.ParseInt(r.PostForm.Get("amount"), 10, 64)
	/*; err != nil {
		res.Status = "error"
		res.Message = fmt.Sprintf("failed getting amount %s from request as int64: %v", r.PostForm.Get("amount"), err)
		wtLogger.Errorf("failed getting amount %s from request as int64: %v", r.PostForm.Get("amount"), err)
	}
	*/
	wtLogger.Debugf("parsed request for url %s: %#v", r.URL.Path, req)

	var t Transfer
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