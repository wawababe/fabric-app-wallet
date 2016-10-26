package task

import "baas/app-wallet/consonlesrvc/common"

var taskLogger = common.NewLogger("task")

type TaskType int
const (
	TASK_TYPE_CREATE_ACCOUNT TaskType = iota
	TASK_TYPE_TRANSFER
)
var TaskTypeMap map[TaskType]string = map[TaskType]string {
	TASK_TYPE_CREATE_ACCOUNT: "createaccount",
	TASK_TYPE_TRANSFER : "accounttransfer",
}


type TaskState int
const (
	TASK_STATE_INIT TaskState = iota
	TASK_STATE_CREATING
	TASK_STATE_TRANSFERING
	TASK_STATE_CHECKING
	TASK_STATE_FIN
	TASK_STATE_FAILED
)

var TaskStateMapCreateAccount map[TaskState]string = map[TaskState]string {
	TASK_STATE_INIT: "pending",
	TASK_STATE_CREATING: "creating",
	TASK_STATE_CHECKING: "checking",
	TASK_STATE_FIN: "fin",
	TASK_STATE_FAILED: "failed",
}

var TaskStateMapTransfer map[TaskState]string = map[TaskState]string {
	TASK_STATE_INIT: "pending",
	TASK_STATE_TRANSFERING: "transferring",
	TASK_STATE_CHECKING: "checking",
	TASK_STATE_FIN: "fin",
	TASK_STATE_FAILED: "failed",
}


