package database

import (
	"testing"
	"database/sql"
	util "baas/app-wallet/consolesrvc/common"
	"encoding/json"
	"fmt"
)

func TestAddTask(t *testing.T) {
	var db *sql.DB
	var err error
	if db, err = sql.Open("mysql", DSN); err != nil {
		dbLogger.Fatal("Failed opening database")
	}

	var testNum = 2

	var useruuid string = "24c78d82-1b28-4dc6-82de-7648ffc87576"
	var accountuuid string = "3a4603e0-a76f-4c77-a9ca-60a717f80e80"
	var account *Account = new(Account)
	if account, err = GetAccount(db, accountuuid); err != nil {
		dbLogger.Errorf("Failed to get account by accountuuid %s and useruuid %s: %v", accountuuid, useruuid, err)
		return
	}

	var transaction *Transaction = new(Transaction)
	var txuuid string = "3684c42b-3e47-45d6-9773-8472c39cef61"
	if transaction, err = GetTransaction(db, txuuid); err != nil {
		dbLogger.Errorf("Failed to get transaction by txuuid %s: %v", txuuid, err)
		return
	}

	var tasks = make([]Task, testNum)
	for i, _ := range tasks {
		tasks[i].TaskUUID = util.GenerateUUID()
		tasks[i].UserUUID = useruuid
		tasks[i].BC_txuuid = util.GenerateUUID()
		if i == 0 {
			tasks[i].Keyword = accountuuid
			payloadBytes, _ := json.Marshal(account)
			tasks[i].Payload = fmt.Sprintf("%s", payloadBytes)
			tasks[i].State = "creating"
			tasks[i].Type = "createaccount"
		}else {
			tasks[i].Keyword = txuuid
			payloadBytes, _ := json.Marshal(transaction)
			tasks[i].Payload = fmt.Sprintf("%s", payloadBytes)
			tasks[i].State = "pending"
			tasks[i].Type = "accounttransfer"
		}
		dbLogger.Debugf("user task: %#v", tasks[i])
	}


	var tests = []struct{
		newline bool
		sep string
		arg *Task
		want int64
	}{
		{false, " ", &tasks[0], 1},
		{true, " ", &tasks[1], 1},
	}

	for i, testitem := range tests {
		rowsAff, _ := AddTask(db, testitem.arg)
		if rowsAff != testitem.want {
			t.Errorf("Test #%d: Add task %#v, affected rows = %d, but want %d", i, testitem.arg, rowsAff, testitem.want)
		}
	}
}

func TestGetTaskByUUID(t *testing.T) {
	var db *sql.DB
	var err error
	if db, err = sql.Open("mysql", DSN); err != nil {
		dbLogger.Fatal("Failed opening database")
	}

//	var useruuid string = "24c78d82-1b28-4dc6-82de-7648ffc87576"
//	var accountuuid string = "3a4603e0-a76f-4c77-a9ca-60a717f80e80"
	var taskuuid string = "10d9ea60-2617-493e-a261-06d46afe480c"

	var task *Task = new(Task)
	if task, err = GetTask(db, taskuuid); err != nil {
		dbLogger.Errorf("failed to get task by taskuuid %s: %v", taskuuid, err)
		t.Errorf("failed to get task by taskuuid %s: %v", taskuuid, err)
	}
	dbLogger.Debugf("succeed in getting task by taskuuid %s: %#v", taskuuid, task)
}

func TestGetTasksByTypeState(t *testing.T) {
	var db *sql.DB
	var err error
	if db, err = sql.Open("mysql", DSN); err != nil {
		dbLogger.Fatal("Failed opening database")
	}

	//	var useruuid string = "24c78d82-1b28-4dc6-82de-7648ffc87576"
	//	var accountuuid string = "3a4603e0-a76f-4c77-a9ca-60a717f80e80"
	//var taskuuid string = "10d9ea60-2617-493e-a261-06d46afe480c"

	var tasks []*Task
	var tasktype = "accounttransfer"
	var taskstate = "pending"
	if tasks, err = GetTasksByTypeState(db, tasktype, taskstate, "creating"); err != nil {
		dbLogger.Errorf("failed to get task by type %s and state %s: %v", tasktype, taskstate, err)
		t.Errorf("failed to get task by type %s and state %s: %v", tasktype, taskstate, err)
	}
	dbLogger.Debugf("succeeded in getting task by type %s and state %s", tasktype, taskstate)

	for i, task := range tasks {
		dbLogger.Debugf("Task #%d: %#v", i, *task)
	}
}

func TestGetTasksByKeywordTypeState(t *testing.T) {
	var db *sql.DB
	var err error
	if db, err = sql.Open("mysql", DSN); err != nil {
		dbLogger.Fatal("Failed opening database")
	}

	//	var useruuid string = "24c78d82-1b28-4dc6-82de-7648ffc87576"
	//	var accountuuid string = "3a4603e0-a76f-4c77-a9ca-60a717f80e80"
	//var taskuuid string = "10d9ea60-2617-493e-a261-06d46afe480c"

	var tasks []*Task
	var tasktype = "accounttransfer"
	var taskstate = "pending"
	var keyword = "3684c42b-3e47-45d6-9773-8472c39cef61"
	if tasks, err = GetTasksByKeywordTypeState(db, keyword, tasktype, taskstate, "creating"); err != nil {
		dbLogger.Errorf("failed to get task by keyword %s, type %s and state %s: %v", keyword, tasktype, taskstate, err)
		t.Errorf("failed to get task by keyword %s, type %s and state %s: %v", keyword, tasktype, taskstate, err)
	}
	dbLogger.Debugf("succeeded in getting task by keyword %s, type %s and state %s", keyword, tasktype, taskstate)

	for i, task := range tasks {
		dbLogger.Debugf("Task #%d: %#v", i, *task)
	}
}

func TestUpdateTaskState(t *testing.T) {
	var db *sql.DB
	var err error
	if db, err = sql.Open("mysql", DSN); err != nil {
		dbLogger.Fatal("Failed opening database")
	}

	//	var useruuid string = "24c78d82-1b28-4dc6-82de-7648ffc87576"
	//	var accountuuid string = "3a4603e0-a76f-4c77-a9ca-60a717f80e80"
	var taskuuid string = "10d9ea60-2617-493e-a261-06d46afe480c"

	var task *Task = new(Task)
	if task, err = GetTask(db, taskuuid); err != nil {
		dbLogger.Errorf("failed to get task by taskuuid %s: %v", taskuuid, err)
		t.Errorf("failed to get task by taskuuid %s: %v", taskuuid, err)
	}
	dbLogger.Debugf("succeed in getting task by taskuuid %s: %#v", taskuuid, task)

	dbLogger.Debugf("task old state is %s", task.State)
	task.State = "fin"
	var rowsAffected int64 = 0
	if rowsAffected, err = UpdateTaskState(db, task); err != nil {
		dbLogger.Errorf("failed to update task state: %v", err)
		t.Errorf("failed to update task state: %v", err)
	}
	dbLogger.Debugf("succeeded in updating task state as %s, affected %d rows", task.State, rowsAffected)


}

func TestUpdateTaskStatePayload(t *testing.T) {

}


