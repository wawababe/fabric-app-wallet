package cronjob

import (
	"os"
	util "baas/app-wallet/consolesrvc/common"
	"database/sql"
	"baas/app-wallet/consolesrvc/database"
	"baas/app-wallet/consolesrvc/wallet/crontask"
	"errors"
	"time"
	"net/http"
	"fmt"
)

var jobLogger = util.NewLogger("CronJob")


type ChaincodeID struct {
	Path string `json:"path,omitempty"`
	Name string `json:"name"`
}
type CtorMsg struct {
	Function string   `json:"function"`
	Args     []string `json:"args"`
}
type Params struct {
	Type          int         `json:"type"`
	ChaincodeID   ChaincodeID `json:"chaincodeID"`
	CtorMsg       CtorMsg     `json:"ctorMsg"`
	SecureContext string      `json:"securecontext, omitempty"`
}

type PeerInvokeReq struct {
	JsonRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	ID      int    `json:"id,omitempty"`
}
type PeerInvokeRes struct {
	JsonRPC string            `json:"jsonrpc"`
	Result  util.BaseResponse `json:"result"`
	ID      int               `json:"id,omitempty"`
}


func NewPeerInvokeReq(funcname string, args []string)(*PeerInvokeReq){
	return &PeerInvokeReq{
		JsonRPC: "2.0",
		Method:  "invoke",
		Params: Params{
			Type: 1,
			ChaincodeID: ChaincodeID{
				Name: "wallet",
			},
			CtorMsg: CtorMsg{
				Function: funcname,
				Args: args,
			},
			//SecureContext: "diego",
		},
		ID: 1,
	}
}


// setTaskState
func setTaskState(db *sql.DB, task *database.Task, newState crontask.TaskState)(affected int64) {
	var err error
	var oldState = task.State
	task.State = newState.String()
	if affected, err = database.UpdateTaskState(db, task); err != nil {
		jobLogger.Errorf("failed updating task %#v from oldstate %s to newstate %s: %v", *task, oldState, task.State, err)
	}
	jobLogger.Debugf("move task state from %s to %s", oldState, newState)
	return affected
}

// todo: this should be augmented by config file
const (
	ENV_DEFAULT_PEER_ADDRESS = "http://127.0.0.1:7050"
)
var ENV_PeerAddress = ""
func MustGetPeerAddress()string{
	ENV_PeerAddress = os.Getenv("PEER_ADDRESS"); //ignore empty PEER_ADDRESS
	if len(ENV_PeerAddress) == 0 {
		jobLogger.Warningf("Environmental variable PEER_ADDRESS not exist, set as default local address: %s", ENV_DEFAULT_PEER_ADDRESS)
		return ENV_DEFAULT_PEER_ADDRESS
	}
	return ENV_PeerAddress
}


// checkBlockchainTxWaitPeer:
// check whether the transaction uuid (the createaccount event response from peer) is valid or not
// POST is retried before DEADLINE with an exponential time back-off
func checkBlockchainTxWaitPeer(db *sql.DB, task *database.Task)(err error){

	var peerAddr string = MustGetPeerAddress()
	const MaxWaitTime = time.Minute * 1
	deadline := time.Now().Add(MaxWaitTime)
	for retries := 0; time.Now().Before(deadline); retries++ {
		queryTransactionResp := func()error{
			var resp *http.Response = new(http.Response)
			resp, err = http.Get(peerAddr + "/" + "transactions" + "/" + task.BC_txuuid);
			if resp.StatusCode == http.StatusOK{
				return nil
			}
			defer resp.Body.Close()
			return errors.New("transaction not found")
		}
		if err = queryTransactionResp(); err == nil {
			break
		}
		time.Sleep(time.Second << uint(retries))
	}

	if err != nil {
		return &JobError{
			OldState: task.State,
			NewState: crontask.STATE_FAILED,
			Err: fmt.Sprintf("failed getting transaction %s from peer %s in %v", task.BC_txuuid, peerAddr, MaxWaitTime),
		}
	}

	return nil
}
