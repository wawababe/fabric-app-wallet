package crontask

import "testing"

func TestAccountCreateTask_Create(t *testing.T) {
	var task CronTask = new(AccountCreateTask)
	var accountuuid = "d769e2dc-2359-4efe-866a-38c4b588fbc5"
	taskuuid, err := task.Create(accountuuid, TYPE_CREATE_ACCOUNT, STATE_INIT)
	if err != nil {
		t.Errorf("failed to create task for accountcreate event: %v", err)
	}
	taskLogger.Debugf("succeeded in creating task %s for accountcreate event", taskuuid)

}

func TestAccountTransferTask_Create(t *testing.T){

}
