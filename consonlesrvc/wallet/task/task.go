package task

type Task interface {
	Create(keyword string, tasktype TaskType, taskstate TaskState)(taskuuid string, err error)
}