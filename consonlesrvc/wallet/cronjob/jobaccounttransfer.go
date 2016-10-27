package cronjob

import (
	util "baas/app-wallet/consonlesrvc/common"
)

var jobLogger = util.NewLogger("JobAccountTransafer")

type JobAccountTransfer struct {

}

func (t *JobAccountTransfer) Run(){


	/*
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
	*/
}
