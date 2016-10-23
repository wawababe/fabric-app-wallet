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
	"os"
	"bytes"
	"strconv"
	"time"
)

type CreateRequest struct {
	authsrvc.AuthRequest
	AccountName string `json:"accountname"`
}

type CreateResponse struct {
	authsrvc.AuthResponse
	AccountUUID string `json:"accountuuid"`
	TxUUID string `json:"txuuid"`
}

type Create struct {
}

func (c *Create) post(req *CreateRequest)(*CreateResponse){
	var res *CreateResponse = new(CreateResponse)
	var err error
	var db *sql.DB = database.GetDB()

	if !req.IsRequestValid(&res.AuthResponse) {
		wtLogger.Warningf("request not valid: %#v", *req)
		res.UserUUID = ""
		return res
	}

	var account *database.Account = new(database.Account)
	if _, err = database.GetAccountByName(db, res.UserUUID + "-" + req.AccountName); err == nil {
		wtLogger.Errorf("failed to create duplicate account %#v", account)
		res.Status = "error"
		res.Message = "failed to create duplicate account"
		return res
	}

	account.UserUUID = res.UserUUID
	account.AccountUUID = util.GenerateUUID()
	account.AccountName = res.UserUUID + "-" + req.AccountName
	account.Amount = 1000
	account.BC_TXUUID = ""
	account.Status = "pending"

	if _, err = database.AddAccount(db, account); err != nil {
		wtLogger.Errorf("failed adding account %#v: %v", account, err)
		res.Status = "error"
		res.Message = "failed adding account"
		return res
	}

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
				Function:"createaccount",
				Args:[]string{""},
			},
			SecureContext:"diego",
		},
		ID: 1,
	}

	peerReq.Params.CtorMsg.Args = []string{
		account.UserUUID,
		account.AccountUUID,
		strconv.FormatInt(int64(account.Amount), 10),
	}
	var peerReqBytes []byte
	if peerReqBytes, err = json.Marshal(peerReq); err != nil {
		wtLogger.Errorf("failed marshalling createAccount request %#v: %v", peerReq, err)
		res.Status = "error"
		res.Message = "failed marshalling createAccount request"
		return res
	}
	wtLogger.Debugf("marshalled createAccount request %#v", string(peerReqBytes))
	var body = new(bytes.Buffer)
	fmt.Fprintf(body, "%s", string(peerReqBytes))
	var peerResp *http.Response = new(http.Response)
	if peerResp, err = http.Post(peerAddr+"/"+"chaincode", "application/json", body); err != nil {
		wtLogger.Errorf("failed posting createAccount request %s to peer %s: %v", string(peerReqBytes), peerAddr, err)
		res.Status = "error"
		res.Message = "failed posint createAccount request to peer"
		return res
	}
	defer peerResp.Body.Close()
	// todo: in this phase, status should be set as creating

	// check whether the createaccount transaction uuid (the response from peer) is valid or not
	// todo: this should be moved out to check
	var invokeResp InvokeRes
	if err = json.NewDecoder(peerResp.Body).Decode(&invokeResp); err != nil {
		wtLogger.Errorf("failed decoding createAccount response from peer: %v", err)
		res.Status = "error"
		res.Message = fmt.Sprintf("failed decoding createAccount response from peer: %v", err)
		return res
	}
	wtLogger.Debugf("decoded createAccount response from peer: %#v", invokeResp)

	res.TxUUID = invokeResp.Result.Message
	if len(res.TxUUID) == 0 {
		res.Status = "error"
		res.Message = "failed creating account, bc_txuuid is empty"
		wtLogger.Error("failed creating account, transaction bc_txuuid is empty")
		return res
	}

	// todo: need to figure out a way to cope with the operation delay of peer
	time.Sleep(time.Second * 4)
	var txreps *http.Response = new(http.Response)
	if txreps, err = http.Get(peerAddr + "/" + "transactions" + "/" + res.TxUUID); err != nil {
		res.Status = "error"
		res.Message = fmt.Sprintf("failed creating account, transaction bc_txuuid %s not exist: %v", res.TxUUID, err)
		wtLogger.Errorf("failed creating account, transaction bc_txuuid %s not exist: %v", res.TxUUID, err)
		return res
	}
 	if txreps.StatusCode != http.StatusOK{
		res.Status = "error"
		res.Message = fmt.Sprintf("failed creating account, transaction bc_txuuid %s not exist", res.TxUUID)
		wtLogger.Errorf("failed creating account, transaction bc_txuuid %s not exist", res.TxUUID)
		return res
	}
	defer txreps.Body.Close()

	account.BC_TXUUID = invokeResp.Result.Message
	account.Status = "created"

	if _, err = database.UpdateAccount(db, account); err != nil {
		wtLogger.Errorf("failed updating account %#v: %v", account, err)
		res.Status = "error"
		res.Message = "failed updating account"
		return res
	}

	res.Status = "ok"
	res.UserUUID = res.UserUUID
	res.Message = "sucessed in creating account"
	res.AccountUUID = account.AccountUUID
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
	fmt.Fprintf(w, "%s", string(resBytes))
}