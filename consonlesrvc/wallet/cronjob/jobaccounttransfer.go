package cronjob

import (
	"database/sql"
	"baas/app-wallet/consonlesrvc/database"
	"baas/app-wallet/consonlesrvc/wallet/crontask"
	"encoding/json"
	"strconv"
	"bytes"
	"fmt"
	"net/http"
)

type JobAccountTransfer struct {

}

func (t *JobAccountTransfer) Run(){
	var tasks []*database.Task
	var err error
	var db *sql.DB = database.GetDB()
	if err = db.Ping(); err != nil {
		jobLogger.Fatal(database.ERROR_DB_NOT_CONNECTED)
		return
	}

	tasks, err = database.GetTasksByTypeState(db, crontask.TypeMap[crontask.TYPE_ACCOUNT_TRANSFER],
		crontask.STATE_INIT.String(),
		crontask.STATE_VALIDATE.String(),
		crontask.STATE_TRANSFER.String(),
		crontask.STATE_CHECK_BCTX.String(),
	)
	if err != nil {
		jobLogger.Warningf("failed fetching task to proccess: %v", err)
		return
	}else if tasks == nil {
		jobLogger.Warning("get zero task to be processed")
		return
	}

	//todo: figure out whether it is a need to process all task concurrently, which may caused a leak of goroutine
	//for task := range tasks {
	//}

	task := tasks[0]
	switch crontask.ParseTaskState(task.State) {
	case crontask.STATE_INIT:
		t.init(db, task)
	case crontask.STATE_VALIDATE:
		t.validate(db, task)
	case crontask.STATE_TRANSFER:
		t.invoke(db, task)
	case crontask.STATE_CHECK_BCTX:
		t.check(db, task)
	}
}


func (t *JobAccountTransfer) init(db *sql.DB, task *database.Task){
	if affected := setTaskState(db, task, crontask.STATE_VALIDATE); affected == 0 {
		return
	}
	jobLogger.Debugf("moved state from %v to %v", crontask.STATE_INIT, crontask.STATE_VALIDATE)
}

// todo: whether it is a need to validate the payload, maybe the payer or the payee not exist any more
func (t *JobAccountTransfer) validate(db *sql.DB, task *database.Task){
	if affected := setTaskState(db, task, crontask.STATE_WAIT_VALIDATE); affected == 0 {
		return
	}
	var err error
	err = validateUserTransaction(db, task);
	if err, ok := err.(*JobError); ok {
		jobLogger.Error(err.Error())
		setTaskState(db, task, err.NewState)
		return
	}
	setTaskState(db, task, crontask.STATE_TRANSFER) //Note: ignore error
}

func (t *JobAccountTransfer) invoke(db *sql.DB, task *database.Task) {
	if affected := setTaskState(db, task, crontask.STATE_WAIT_TRANSFER); affected == 0 {
		return
	}
	var err error
	err = transfer(db, task)
	if err, ok := err.(*JobError); ok {
		jobLogger.Error(err.Error())
		setTaskState(db, task, err.NewState)
		return
	}
	setTaskState(db, task, crontask.STATE_CHECK_BCTX) //Note: ignore error
}

func (t *JobAccountTransfer) check(db *sql.DB, task *database.Task){
	if affected := setTaskState(db, task, crontask.STATE_WAIT_CHECK_BCTX); affected == 0 {
		return
	}
	var err error
	err = checkBlockchainTxWaitPeer(db, task)
	if err, ok := err.(*JobError); ok {
		jobLogger.Error(err.Error())
		setTaskState(db, task, err.NewState)
		return
	}
	if err == nil {
		var tx = new(database.Transaction)
		json.Unmarshal([]byte(task.Payload), tx) //note: ignore errors
		var payer = new(database.Account)
		var payee = new(database.Account)

		//todo: CAREAT: this should never happen, could be prevented from validate phase
		//todo: also, account should not be deleted if it is referenced by a non-failed task
		if payer, err = database.GetAccountByAccountID(db, tx.PayerAccountID); err != nil {
			jobLogger.Fatalf("failed to get payer account by accountid %s: %v", tx.PayerAccountID, err)
		}
		if payee, err = database.GetAccountByAccountID(db, tx.PayeeAccountID); err != nil {
			jobLogger.Fatalf("failed to get payee account by accountid %S: %v", tx.PayeeAccountID, err)
		}
		// todo: this should be done with further consideration, two mode:
		// 1. firstly update account, once failed in blockchain, then recover the account;
		// 2. or just try to transfer in blockchain, then query blockchain to update the account
		payer.Amount -= tx.Amount
		payee.Amount += tx.Amount
		if _, err = database.UpdateAccount(db, payer); err != nil {
			jobLogger.Fatalf("failed to update payer account as %#v: %v", *payer, err)
		}
		if _, err = database.UpdateAccount(db, payee); err != nil {
			jobLogger.Fatalf("failed to update payee account as %#v: %v", *payee, err)
		}
		tx.BC_txuuid = task.BC_txuuid
		tx.Status = "transferred"
		if _, err = database.UpdateTransaction(db, tx); err != nil {
			jobLogger.Fatalf("failed to update transaction as %#v: %v", *tx, err)
		}

	}
	setTaskState(db, task, crontask.STATE_FIN) //Note: ignore error
}



// todo: extract the common behavior between different jobs;
func validateUserTransaction(db *sql.DB, task *database.Task)(err error){
	if err = db.Ping(); err != nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.ParseTaskState(task.State),
			Err: database.ERROR_DB_NOT_CONNECTED,
		}
	}

	var user *database.User
	if user, err = database.GetUser(db, task.UserUUID); err != nil || user == nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("user %s in the task not exist: %v", task.UserUUID, err),
		}
	}

	var tx *database.Transaction
	if tx, err = database.GetTransaction(db, task.Keyword); err != nil || tx == nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("transaction %s in the task not exist: %v", task.Keyword, err),
		}
	}

	var txpayload *database.Transaction = new(database.Transaction)
	if err = json.Unmarshal([]byte(task.Payload), txpayload); err != nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("failed unmarshalling transaction payload %s: %v", task.Payload, err),
		}
	}

	var payer = new(database.Account)
	var payee = new(database.Account)
	if payer, err = database.GetAccountByAccountID(db, txpayload.PayerAccountID); err != nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("failed getting payer account by accountid %s: %v", txpayload.PayerAccountID, err),
		}
	}
	if payee, err = database.GetAccountByAccountID(db, txpayload.PayeeAccountID); err != nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("failed getting payee account by accountid %s: %v", txpayload.PayeeAccountID, err),
		}
	}

	if payer == nil || payee == nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("payer account %s or payee account %s not exist", payer.AccountID, payee.AccountID),
		}
	}

	if payer.UserUUID != task.UserUUID {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("user %s is not the owner of the account %s", task.UserUUID, txpayload.PayerAccountID),
		}
	}
	return nil
}



func transfer(db *sql.DB, task *database.Task)(err error){
	var peerAddr string = MustGetPeerAddress()

	var payload *database.Transaction = new(database.Transaction)
	if err = json.Unmarshal([]byte(task.Payload), payload); err != nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("failed unmarshalling task payload %s as transaction: %v", task.Payload, err),
		}
	}

	var payer *database.Account = new(database.Account)
	var payee *database.Account = new(database.Account)

	payer, _ = database.GetAccountByAccountID(db, payload.PayerAccountID) //Note: ignore error, assumed that it has passed the validate phase
	payee, _ = database.GetAccountByAccountID(db, payload.PayeeAccountID) //Note: ignore error, assumed that it has passed the validate phase
	if payer == nil || payee == nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("payer %#v or payee %#v not exist", *payer, *payee),
		}
	}

	var args []string = []string{
		task.UserUUID,
		payer.AccountUUID,
		payee.AccountUUID,
		strconv.FormatInt(int64(payload.Amount), 10),
	}
	var peerReq = NewPeerInvokeReq("accounttransfer", args)
	var peerReqBytes []byte
	if peerReqBytes, err = json.Marshal(peerReq); err != nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("failed marshalling accounttransfer request %#v: %v", peerReq, err),
		}
	}
	jobLogger.Debugf("marshalled accounttransfer request %#v", string(peerReqBytes))

	var body = new(bytes.Buffer)
	fmt.Fprintf(body, "%s", string(peerReqBytes))
	var peerResp *http.Response = new(http.Response)


	if peerResp, err = http.Post(peerAddr+"/"+"chaincode", "application/json", body); err != nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.ParseTaskState(task.State),
			Err: fmt.Sprintf("failed posting accounttransfre request %s to peer %s: %v", string(peerReqBytes), peerAddr, err),
		}
	}
	defer peerResp.Body.Close()

	// todo: in this phase, status should be set as creating
	var invokeResp PeerInvokeRes
	if err = json.NewDecoder(peerResp.Body).Decode(&invokeResp); err != nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("failed decoding accounttransfer responsonse from peer: %v", err),
		}

	}
	jobLogger.Debugf("decoded createAccount response from peer: %#v", invokeResp)

	task.BC_txuuid = invokeResp.Result.Message
	if len(task.BC_txuuid) == 0 {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprint("failed creating account, transaction bc_txuuid is empty"),
		}
	}

	if _, err := database.UpdateTask(db, task); err != nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("failed adding blockchain txuuid %s into task %#v: %v", task.BC_txuuid, *task, err),
		}
	}
	jobLogger.Debugf("update task as %#v", *task)
	return nil
}

