package crontask

import (
	"baas/app-wallet/consolesrvc/common"
)

var taskLogger = common.NewLogger("task")

type TaskType int
const (
	TYPE_CREATE_ACCOUNT TaskType = iota
	TYPE_ACCOUNT_TRANSFER
)
var TypeMap map[TaskType]string = map[TaskType]string {
	TYPE_CREATE_ACCOUNT: "createaccount",
	TYPE_ACCOUNT_TRANSFER : "accounttransfer",
}


type TaskState int
const (
	STATE_INIT TaskState = iota
	STATE_VALIDATE //validate whether the task is valid or not
	STATE_WAIT_VALIDATE

	STATE_CREATE_ACCOUNT
	STATE_WAIT_CREATE_ACCOUNT
	STATE_TRANSFER
	STATE_WAIT_TRANSFER

	// check block chain transaction by
	// employing a strategy which retries after exponential time in finite number of steps when failed
	STATE_CHECK_BCTX
	STATE_WAIT_CHECK_BCTX

	STATE_FIN
	STATE_FAILED
)


func (ts TaskState) String()string {
	var stateMap map[TaskState]string = map[TaskState]string{
		STATE_INIT: "pending",
		STATE_VALIDATE: "validate",
		STATE_WAIT_VALIDATE: "wait_validate",
		STATE_CREATE_ACCOUNT: "create_account",
		STATE_WAIT_CREATE_ACCOUNT: "wait_create_account",
		STATE_TRANSFER: "transfer",
		STATE_WAIT_TRANSFER: "wait_transfer",
		STATE_CHECK_BCTX: "check_bctx",
		STATE_WAIT_CHECK_BCTX: "wait_check_bctx",
		STATE_FIN: "fin",
		STATE_FAILED: "failed",
	}
	return stateMap[ts]
}

func ParseTaskState(ts string)TaskState {
	var stateMap map[string]TaskState = map[string]TaskState{
		"pending": STATE_INIT,
		"validate": STATE_VALIDATE,
		"wait_validate": STATE_WAIT_VALIDATE,
		"create_account": STATE_CREATE_ACCOUNT,
		"wait_create_account": STATE_WAIT_CREATE_ACCOUNT,
		"transfer": STATE_TRANSFER,
		"wait_transfer": STATE_WAIT_TRANSFER,
		"check_bctx": STATE_CHECK_BCTX,
		"wait_check_bctx": STATE_WAIT_CHECK_BCTX,
		"fin": STATE_FIN,
		"failed": STATE_FAILED,
	}
	return stateMap[ts]
}



