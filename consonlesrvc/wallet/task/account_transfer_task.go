package task

import (
	"baas/app-wallet/consonlesrvc/database"
	"errors"
	"encoding/json"
	"fmt"
	util "baas/app-wallet/consonlesrvc/common"
)


type AccountTransferTask struct{

}

func (t *AccountTransferTask) Create(txuuid string, tasktype TaskType, taskstate TaskState)(taskuuid string, err error){
	var db = database.GetDB()
	if err = db.Ping(); err != nil {
		taskLogger.Fatal(database.ERROR_DB_NOT_CONNECTED)
		return "", errors.New(database.ERROR_DB_NOT_CONNECTED)
	}
	var tx *database.Transaction = new(database.Transaction)
	if tx, err = database.GetTransaction(db, txuuid); err != nil {
		taskLogger.Errorf("failed to get transaction by txuuid %s", txuuid)
		return "", fmt.Errorf("failed to get transaction by txuuid %s", txuuid)
	}

	var task *database.Task = new(database.Task)
	task.TaskUUID = util.GenerateUUID()
	task.UserUUID = tx.PayerUUID
	task.Keyword = txuuid
	task.BC_txuuid = ""
	task.Type = TaskTypeMap[tasktype]
	task.State = TaskStateMapTransfer[taskstate]
	payloadBytes, _ := json.Marshal(tx)
	task.Payload = string(payloadBytes)

	if _, err = database.AddTask(db, task); err != nil {
		taskLogger.Errorf("failed to add task %#v into database: %v", *task, err)
		return "", fmt.Errorf("failed to add task %#v into database: %v", *task, err)
	}
	taskLogger.Debugf("success in adding task into database:\n%#v", *task)
	taskuuid = task.TaskUUID
	return taskuuid, nil
}
