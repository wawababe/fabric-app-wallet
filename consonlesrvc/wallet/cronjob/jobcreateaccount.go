package cronjob

import (
	"baas/app-wallet/consonlesrvc/database"
	"database/sql"
	util "baas/app-wallet/consonlesrvc/common"
	"baas/app-wallet/consonlesrvc/wallet/crontask"
	"fmt"
	"os"
	"bytes"
	"encoding/json"
	"strconv"
	"net/http"
)

var jobLogger = util.NewLogger("JobCreateAccount")

type JobCreateAccount struct {
}

type JobError struct {
	OldState string
	NewState string
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
		crontask.StateMap[crontask.STATE_INIT],
		crontask.StateMap[crontask.STATE_VALIDATE],
		crontask.StateMap[crontask.STATE_CREATE_ACCOUNT],
		crontask.StateMap[crontask.STATE_CHECK_BCTX],
	)
	if err != nil || tasks == nil {
		jobLogger.Warning("failed fetching task to proccess: %v", err)
		return
	}

	//todo: figure out whether it is a need to process all task concurrently, which may caused a leak of goroutine
	//for task := range tasks {
	//}

	task := tasks[0]
	switch task.State {
	case crontask.StateMap[crontask.STATE_INIT]:
		t.processPending(db, task)
	case crontask.StateMap[crontask.STATE_VALIDATE]:
		t.validate(db, task)
	case crontask.StateMap[crontask.STATE_CREATE_ACCOUNT]:
	}




	/*

	// todo: need to figure out a way to cope with the operation delay of peer
	time.Sleep(time.Second * 3)
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
*/
}

func (t *JobCreateAccount) processPending(db *sql.DB, task *database.Task){
	if affected := setTaskState(db, task, crontask.StateMap[crontask.STATE_VALIDATE]); affected == 0 {
		return
	}
	jobLogger.Debug("moved state from pending to %s", crontask.StateMap[crontask.STATE_VALIDATE])
}
func (t *JobCreateAccount) validate(db *sql.DB, task *database.Task){
	if affected := setTaskState(db, task, crontask.StateMap[crontask.STATE_WAIT_VALIDATE]); affected == 0 {
		return
	}
	var err error
	err = validateUserAccount(db, task);
	if err, ok := err.(JobError); ok {
		jobLogger.Error(err.Error())
		setTaskState(db, task, err.NewState)
	}
	setTaskState(db, task, crontask.StateMap[crontask.STATE_CREATE_ACCOUNT])
}

func (t *JobCreateAccount) createAccount(db *sql.DB, task *database.Task){
	if affected := setTaskState(db, task, crontask.StateMap[crontask.STATE_WAIT_CREATE_ACCOUNT]); affected == 0 {
		return
	}
	var err error
	err = createAccountOnBlockchain(db, task)
	if err != nil {
		//todo
	}
}

func (t *JobCreateAccount) checkBlockchainTX(db *sql.DB, task *database.Task){
	if affected := setTaskState(db, task, crontask.StateMap[crontask.STATE_WAIT_CHECK_BCTX]); affected == 0 {
		return
	}
}


func setTaskState(db *sql.DB, task *database.Task, newstate string)(affected int64) {
	var err error
	var oldstate = task.State
	task.State = newstate
	if affected, err = database.UpdateTaskState(db, task); err != nil {
		jobLogger.Errorf("failed updating task %#v from oldstate %s to newstate %s: %v", *task, oldstate, newstate, err)
	}
	jobLogger.Debugf("move task state from %s to %s", oldstate, newstate)
	return affected
}

// validate whether the user has the account
func validateUserAccount(db *sql.DB, task *database.Task)(err error){
	if err = db.Ping(); err != nil {
		return &JobError{
			OldState: task.State,
			NewState: task.State,
			Err: database.ERROR_DB_NOT_CONNECTED,
		}
	}

	var user *database.User
	if user, err = database.GetUser(db, task.UserUUID); err != nil || user == nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.StateMap[crontask.STATE_FAILED],
			Err: fmt.Sprintf("user %s in the task not exist: %v", task.UserUUID, err),
		}
	}

	var account *database.Account
	if account, err = database.GetAccount(db, task.Keyword); err != nil || account == nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.StateMap[crontask.STATE_FAILED],
			Err: fmt.Errorf("account %s in the task not exist: %v", task.Keyword, err),
		}
	}
	return nil
}

func createAccountOnBlockchain(db *sql.DB, task *database.Task)(err error){
	var peerAddr string
	if peerAddr = os.Getenv("PEER_ADDRESS"); len(peerAddr) == 0 {
		jobLogger.Fatal("failed getting environmental variable PEER_ADDRESS")
		return &JobError{
			OldState: task.State,
			NewState: task.State,
			Err: "failed getting environmental variable PEER_ADDRESS",
		}
	}

	var account *database.Account = new(database.Account)
	if err = json.Unmarshal([]byte(task.Payload), account); err != nil {
		jobLogger.Fatalf("failed unmarshalling task payload %s as account: %v", task.Payload, err)
		return &JobError{
			OldState: task.State,
			NewState: crontask.StateMap[crontask.STATE_FAILED],
			Err: fmt.Errorf("failed unmarshalling task payload %s as account: %v", task.Payload, err),
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
		jobLogger.Errorf("failed marshalling createAccount request %#v: %v", peerReq, err)
		return &JobError{
			OldState: task.State,
			NewState: crontask.StateMap[crontask.STATE_FAILED],
			Err: fmt.Errorf("failed marshalling createAccount request %#v: %v", peerReq, err),
		}
	}
	jobLogger.Debugf("marshalled createAccount request %#v", string(peerReqBytes))

	var body = new(bytes.Buffer)
	fmt.Fprintf(body, "%s", string(peerReqBytes))
	var peerResp *http.Response = new(http.Response)
	if peerResp, err = http.Post(peerAddr+"/"+"chaincode", "application/json", body); err != nil {
		jobLogger.Errorf("failed posting createAccount request %s to peer %s: %v", string(peerReqBytes), peerAddr, err)
		return &JobError{
			OldState: task.State,
			NewState: task.State,
			Err: fmt.Errorf("failed marshalling createAccount request %#v: %v", peerReq, err),
		}
	}
	defer peerResp.Body.Close()
	// todo: in this phase, status should be set as creating

	// check whether the createaccount transaction uuid (the response from peer) is valid or not
	// todo: this should be moved out to check
	var invokeResp PeerInvokeRes
	if err = json.NewDecoder(peerResp.Body).Decode(&invokeResp); err != nil {
		jobLogger.Errorf("failed decoding createAccount response from peer: %v", err)

		return
	}
	jobLogger.Debugf("decoded createAccount response from peer: %#v", invokeResp)

	task.BC_txuuid = invokeResp.Result.Message
	if len(task.BC_txuuid) == 0 {
		jobLogger.Error("failed creating account, transaction bc_txuuid is empty")
		return
	}
}

