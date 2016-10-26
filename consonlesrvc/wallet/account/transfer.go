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
	"bytes"
	"os"
	"time"
	"strings"
	util "baas/app-wallet/consonlesrvc/common"
)

type TransferRequest struct {
	authsrvc.AuthRequest
	PayerAccountID string `json:"payeraccountid"`
	PayeeAccountID string `json:"payeeaccountid"`
	Amount int64
}

type TransferResponse struct {
	authsrvc.AuthResponse
	TxUUID string `json:"txuuid,omitempty"`
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

	// todo: should not check, this should be done by chaincode
	if payerAccount.Amount < req.Amount {
		wtLogger.Errorf("failed to pay, payer's account only has %d, not enough to pay %d", payerAccount.Amount, req.Amount)
		res.Status = "error"
		res.Message = fmt.Sprintf("failed to pay, payer's account only has %d, not enough to pay %d", payerAccount.Amount, req.Amount)
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

	// todo: add a task for account transfer

	// todo: run these in a goroutine, which could be implemented by cron mechanism and additional task table
	var peerAddr string
	if peerAddr = os.Getenv("PEER_ADDRESS"); len(peerAddr) == 0 {
		wtLogger.Fatal("failed getting environmental variable PEER_ADDRESS")
		res.Status = "error"
		res.Message = "failed getting peer address"
		return res
	}

	var peerReq = PeerReq{
		JsonRPC: "2.0",
		Method: "invoke",
		Params: Params{
			Type: 1,
			ChaincodeID: ChaincodeID{
				Name:"wallet",
			},
			CtorMsg: CtorMsg{
				Function:"accounttransfer",
				Args:[]string{""},
			},
			SecureContext:"diego",
		},
		ID: 1,
	}

	peerReq.Params.CtorMsg.Args = []string{
		tx.TxUUID,
		payerAccount.AccountUUID,
		payeeAccount.AccountUUID,
		strconv.FormatInt(int64(tx.Amount), 10),
	}
	var peerReqBytes []byte
	if peerReqBytes, err = json.Marshal(peerReq); err != nil {
		wtLogger.Errorf("failed marshalling account transfer request %#v: %v", peerReq, err)
		res.Status = "error"
		res.Message = "failed marshalling account transfer request"
		return res
	}
	wtLogger.Debugf("marshalled createAccount request %#v", string(peerReqBytes))
	var body = new(bytes.Buffer)
	fmt.Fprintf(body, "%s", string(peerReqBytes))

	var peerResp *http.Response = new(http.Response)
	if peerResp, err = http.Post(peerAddr+"/"+"chaincode", "application/json", body); err != nil {
		wtLogger.Errorf("failed posting account transfer request %s to peer %s: %v", string(peerReqBytes), peerAddr, err)
		res.Status = "error"
		res.Message = "failed posint account transfer request to peer"
		return res
	}
	defer peerResp.Body.Close()
	// todo: in this phase, status should be set as transferring

	// check whether the transferaccount transaction uuid (the response from peer) is valid or not
	// todo: this should be moved out to check
	var invokeResp InvokeRes
	if err = json.NewDecoder(peerResp.Body).Decode(&invokeResp); err != nil {
		wtLogger.Errorf("failed decoding account transfer response from peer: %v", err)
		res.Status = "error"
		res.Message = fmt.Sprintf("failed decoding account transfer response from peer: %v", err)
		return res
	}
	wtLogger.Debugf("decoded account transfer response from peer: %#v", invokeResp)

	res.TxUUID = invokeResp.Result.Message
	if len(res.TxUUID) == 0 {
		res.Status = "error"
		res.Message = "failed transferring account, bc_txuuid is empty"
		wtLogger.Error("failed transferring account, transaction bc_txuuid is empty")
		return res
	}

	// todo: need to figure out a way to cope with the operation delay of peer
	time.Sleep(time.Second * 4)
	var txreps *http.Response = new(http.Response)
	if txreps, err = http.Get(peerAddr + "/" + "transactions" + "/" + res.TxUUID); err != nil {
		res.Status = "error"
		res.Message = fmt.Sprintf("failed transferring account, transaction bc_txuuid %s not exist: %v", res.TxUUID, err)
		wtLogger.Errorf("failed transferring account, transaction bc_txuuid %s not exist: %v", res.TxUUID, err)
		return res
	}
	if txreps.StatusCode != http.StatusOK{
		res.Status = "error"
		res.Message = fmt.Sprintf("failed transferring account, transaction bc_txuuid %s not exist", res.TxUUID)
		wtLogger.Errorf("failed transferring account, transaction bc_txuuid %s not exist", res.TxUUID)
		return res
	}
	defer txreps.Body.Close()

	tx.BC_txuuid = invokeResp.Result.Message
	tx.Status = "created"

	if _, err = database.UpdateTransaction(db, tx); err != nil {
		wtLogger.Errorf("failed updating transaction %#v: %v", tx, err)
		res.Status = "error"
		res.Message = "failed updating transaction"
		return res
	}

	payerAccount.Amount -= req.Amount
	payeeAccount.Amount += req.Amount
	if _, err = database.UpdateAccount(db, payerAccount); err != nil {
		wtLogger.Errorf("failed updating payer account %#v: %v", payerAccount, err)
		res.Status = "error"
		res.Message = "failed updating payer account"
		return res
	}

	if _, err = database.UpdateAccount(db, payeeAccount); err != nil {
		wtLogger.Errorf("failed updating payee account %#v: %v", payerAccount, err)
		res.Status = "error"
		res.Message = "failed updating payee account"
		return res
	}

	res.Status = "ok"
	res.UserUUID = res.UserUUID
	res.Message = "sucessed in creating account"
	res.TxUUID = tx.TxUUID
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