package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type Task struct {
	RowID int64
	TaskUUID string
	UserUUID string
	Keyword string
	BC_txuuid string
	Type string
	State string
	Payload string
}


func AddTask(db *sql.DB, t *Task)(int64, error){
	dbLogger.Debug("AddTask...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("INSERT INTO task(taskuuid, useruuid, keyword, bc_txuuid, type, state, payload) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	addResult, err = stmt.Exec(t.TaskUUID, t.UserUUID, t.Keyword, t.BC_txuuid, t.Type, t.State, t.Payload)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
}


// GetTask: get task by taskuuid
func GetTask(db *sql.DB, taskuuid string) (*Task, error) {
	dbLogger.Debug("GetTask...")
	var task = new(Task)
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return nil, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("SELECT rowid, taskuuid, useruuid, keyword, bc_txuuid, type, state, payload FROM task WHERE taskuuid = ? and deleted = 0")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return nil, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(taskuuid).Scan(&task.RowID, &task.TaskUUID, &task.UserUUID, &task.Keyword, &task.BC_txuuid, &task.Type, &task.State, &task.Payload); err != nil {
		dbLogger.Errorf("Failed getting task by taskuuid %s: %v", taskuuid, err)
		return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
	}
	dbLogger.Debugf("Get task by taskuuid %s: \n%#v", taskuuid, *task)

	return task, nil
}

// get tasks by task type and task state
func GetTasksByTypeState(db *sql.DB, tasktype string, states ...string)([]*Task, error){
	dbLogger.Debug("GetTasksByTypeState...")
	var err error
	var rows *sql.Rows
	var tasks []*Task

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return nil, errors.New(ERROR_DB_NOT_CONNECTED)
	}


	for _, state := range states {
		if rows, err = db.Query("SELECT rowid, taskuuid, useruuid, keyword, bc_txuuid, type, state, payload FROM task where type = ? and state = ? and deleted = 0", &tasktype, &state); err != nil {
			dbLogger.Errorf("Failed getting tasks by type %s and state %v: %v", tasktype, state, err)
			continue
			//return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
		}

		for rows.Next() {
			var t *Task = new(Task)
			if err := rows.Scan(&t.RowID, &t.TaskUUID, &t.UserUUID, &t.Keyword, &t.BC_txuuid, &t.Type, &t.State, &t.Payload); err != nil {
				dbLogger.Fatal(err)
				continue //todo: caution: error should not happen here, any potential error should be considered as data type error
			}
			dbLogger.Debugf("task with type %s and states %v: %#v", tasktype, states, *t)
			tasks = append(tasks, t)
		}
	}
	defer rows.Close()

	return tasks, nil
}


// get tasks by keyword, task type and task state
func GetTasksByKeywordTypeState(db *sql.DB, keyword string, tasktype string, states ...string)([]*Task, error){
	dbLogger.Debug("GetTasksByKeywordTypeState...")
	var err error
	var rows *sql.Rows
	var tasks []*Task

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return nil, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	for _, state := range states {
		if rows, err = db.Query("SELECT rowid, taskuuid, useruuid, keyword, bc_txuuid, type, state, payload FROM task where keyword = ? and type = ? and state = ? and deleted = 0", &keyword, &tasktype, &state); err != nil {
			dbLogger.Errorf("Failed getting tasks by keyword %s and type %s and state %s: %v", keyword, tasktype, state, err)
			continue //todo: error should be further considered, if it is caused by non-exist , then it should be ignored here
			//return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
		}

		for rows.Next() {
			var t *Task = new(Task)
			if err := rows.Scan(&t.RowID, &t.TaskUUID, &t.UserUUID, &t.Keyword, &t.BC_txuuid, &t.Type, &t.State, &t.Payload); err != nil {
				dbLogger.Fatal(err)
				continue //just ignore scan error
			}
			fmt.Println("Task states:", states)
			dbLogger.Debugf("task with type %s: %#v", tasktype, *t)
			tasks = append(tasks, t)
		}
	}
	defer rows.Close()
	//todo: should return error if tasks is nil
	return tasks, nil
}

func UpdateTaskState(db *sql.DB, t *Task)(int64, error){
	dbLogger.Debug("UpdateTaskState...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("UPDATE task SET state = ? WHERE taskuuid = ? and deleted = 0")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	addResult, err = stmt.Exec(t.State, t.TaskUUID)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
}


func UpdateTaskStatePayload(db *sql.DB, t *Task)(int64, error){
	dbLogger.Debug("UpdateTaskStatePayload...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("UPDATE task SET state = ?, payload = ? WHERE taskuuid = ? and deleted = 0")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	addResult, err = stmt.Exec(t.State, t.Payload, t.TaskUUID)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
}

func UpdateTask(db *sql.DB, t *Task)(int64, error){
	dbLogger.Debug("UpdateTask...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("UPDATE task SET bc_txuuid = ?, state = ?, payload = ? WHERE taskuuid = ? and deleted = 0")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()
	dbLogger.Debugf("expecting to update task: %#v", *t)
	addResult, err = stmt.Exec(t.BC_txuuid, t.State, t.Payload, t.TaskUUID)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
}

func DeleteTask(db *sql.DB, t *Task)(int64, error){
	dbLogger.Debug("DeleteTask...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("UPDATE task SET deleted = 1 WHERE taskuuid = ?")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()
	dbLogger.Debugf("expecting to delete task: %#v", *t)
	addResult, err = stmt.Exec(t.TaskUUID)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
}