package crontask

type CronTask interface {
	Create(keyword string, tasktype TaskType, taskstate TaskState)(taskuuid string, err error)
}