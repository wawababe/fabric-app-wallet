package crontask

import (
	"baas/app-wallet/consonlesrvc/database"
	"fmt"
	util "baas/app-wallet/consonlesrvc/common"
	"encoding/json"
	"errors"
)

type AccountCreateTask struct {

}

func (t *AccountCreateTask) Create(accountuuid string, tasktype TaskType, taskstate TaskState)(taskuuid string, err error){
	var db = database.GetDB()
	if err = db.Ping(); err != nil {
		taskLogger.Fatal(database.ERROR_DB_NOT_CONNECTED)
		return "", errors.New(database.ERROR_DB_NOT_CONNECTED)
	}
	var account *database.Account = new(database.Account)
	if account, err = database.GetAccount(db, accountuuid); err != nil {
		taskLogger.Errorf("failed to get account by accountuuid %s", accountuuid)
		return "", fmt.Errorf("failed to get account by accountuuid %s", accountuuid)
	}

	var task *database.Task = new(database.Task)
	task.TaskUUID = util.GenerateUUID()
	task.UserUUID = account.UserUUID
	task.Keyword = accountuuid
	task.BC_txuuid = ""
	task.Type = TypeMap[tasktype]
	task.State = taskstate.String()
	payloadBytes, _ := json.Marshal(account)
	task.Payload = string(payloadBytes)

	if _, err = database.AddTask(db, task); err != nil {
		taskLogger.Errorf("failed to add task %#v into database: %v", *task, err)
		return "", fmt.Errorf("failed to add task %#v into database: %v", *task, err)
	}
	taskLogger.Debugf("success in adding task into database:\n%#v", *task)
	taskuuid = task.TaskUUID
	return taskuuid, nil
}
