package cronjob

import (
	"baas/app-wallet/consolesrvc/database"
	"database/sql"
	"baas/app-wallet/consolesrvc/wallet/crontask"
	"fmt"
	"bytes"
	"encoding/json"
	"strconv"
	"net/http"

)

type JobCreateAccount struct {
}

type JobError struct {
	OldState string
	NewState crontask.TaskState
	Err string
}
func (t *JobError)Error()string{
	return t.Err
}

func (t *JobCreateAccount) Run(){

	var tasks []*database.Task
	var err error
	var db *sql.DB = database.GetDB()
	if err = db.Ping(); err != nil {
		jobLogger.Fatal(database.ERROR_DB_NOT_CONNECTED)
		return
	}

	tasks, err = database.GetTasksByTypeState(db, crontask.TypeMap[crontask.TYPE_CREATE_ACCOUNT],
		crontask.STATE_INIT.String(),
		crontask.STATE_VALIDATE.String(),
		crontask.STATE_CREATE_ACCOUNT.String(),
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
	case crontask.STATE_CREATE_ACCOUNT:
		t.invoke(db, task)
	case crontask.STATE_CHECK_BCTX:
		t.check(db, task)
	}


}

func (t *JobCreateAccount) init(db *sql.DB, task *database.Task){
	if affected := setTaskState(db, task, crontask.STATE_VALIDATE); affected == 0 {
		return
	}
	jobLogger.Debugf("moved state from pending to %v", crontask.STATE_VALIDATE)
}
func (t *JobCreateAccount) validate(db *sql.DB, task *database.Task){
	if affected := setTaskState(db, task, crontask.STATE_WAIT_VALIDATE); affected == 0 {
		return
	}
	var err error
	err = validateUserAccount(db, task);
	if err, ok := err.(*JobError); ok {
		jobLogger.Error(err.Error())
		setTaskState(db, task, err.NewState)
		return
	}
	setTaskState(db, task, crontask.STATE_CREATE_ACCOUNT) //Note: ignore error
}

func (t *JobCreateAccount) invoke(db *sql.DB, task *database.Task) {
	if affected := setTaskState(db, task, crontask.STATE_WAIT_CREATE_ACCOUNT); affected == 0 {
		return
	}
	var err error
	err = createAccount(db, task)
	if err, ok := err.(*JobError); ok {
		jobLogger.Error(err.Error())
		setTaskState(db, task, err.NewState)
		return
	}
	setTaskState(db, task, crontask.STATE_CHECK_BCTX) //Note: ignore error
}

func (t *JobCreateAccount) check(db *sql.DB, task *database.Task){
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
	if err == nil { //success to create account on blockchain, update accounts in local database
		var account *database.Account = new(database.Account)
		json.Unmarshal([]byte(task.Payload), account); //Note: ignore error
		account.BC_TXUUID = task.BC_txuuid
		account.Status = "created"
		if _, err = database.UpdateAccount(db, account); err != nil { //Nore: ignore error
			jobLogger.Errorf("failed updating account %#v: %v", account, err)
		}
	}
	setTaskState(db, task, crontask.STATE_FIN) //Note: ignore error
}


// validateUserAccount: validate whether the user has the account
// todo: extract the common behavior between different jobs;
// todo: validate whether the user has the account
func validateUserAccount(db *sql.DB, task *database.Task)(err error){
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

	var account *database.Account
	if account, err = database.GetAccount(db, task.Keyword); err != nil || account == nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("account %s in the task not exist: %v", task.Keyword, err),
		}
	}
	return nil
}

func createAccount(db *sql.DB, task *database.Task)(err error){
	var peerAddr string = MustGetPeerAddress()

	var account *database.Account = new(database.Account)
	if err = json.Unmarshal([]byte(task.Payload), account); err != nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("failed unmarshalling task payload %s as account: %v", task.Payload, err),
		}
	}

	var args []string = []string{
		account.UserUUID,
		account.AccountUUID,
		strconv.FormatInt(int64(account.Amount), 10),
	}
	var peerReq = NewPeerInvokeReq("createaccount", args)
	var peerReqBytes []byte
	if peerReqBytes, err = json.Marshal(peerReq); err != nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("failed marshalling createAccount request %#v: %v", peerReq, err),
		}
	}
	jobLogger.Debugf("marshalled createAccount request %#v", string(peerReqBytes))

	var body = new(bytes.Buffer)
	fmt.Fprintf(body, "%s", string(peerReqBytes))
	var peerResp *http.Response = new(http.Response)


	if peerResp, err = http.Post(peerAddr+"/"+"chaincode", "application/json", body); err != nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.ParseTaskState(task.State),
			Err: fmt.Sprintf("failed posting createAccount request %s to peer %s: %v", string(peerReqBytes), peerAddr, err),
		}
	}
	defer peerResp.Body.Close()

	// todo: in this phase, status should be set as creating
	var invokeResp PeerInvokeRes
	if err = json.NewDecoder(peerResp.Body).Decode(&invokeResp); err != nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("failed decoding createAccount responsonse from peer: %v", err),
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

