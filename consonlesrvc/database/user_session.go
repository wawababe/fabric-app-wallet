package database

import (
	"time"
	"database/sql"
	"errors"
	"fmt"
)

type UserSession struct {
	RowID int64
	UserUUID string
	SessionUUID string
	ExpiredAt string
}

func (t *UserSession) AddExpiredTimeByDays(days int){
	if  len(t.ExpiredAt) == 0 {
		t.ExpiredAt = time.Now().Format(DATETIME_FORMAT)
	}
	tm, _ := time.Parse(DATETIME_FORMAT,t.ExpiredAt)
	tm.AddDate(0, 0, days)
	t.ExpiredAt = tm.Format(DATETIME_FORMAT)

}


// AddUserSession: insert a new user session into table userSession
func AddUserSession(db *sql.DB, us *UserSession)(int64, error) {
	dbLogger.Debug("AddUserSession...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("INSERT INTO USERSESSION(useruuid, sessionuuid, expiredAt) VALUES(?, ?, ?)")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	addResult, err = stmt.Exec(us.UserUUID, us.SessionUUID, []byte(us.ExpiredAt))
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
}

// GetUserSession: query userSession table by userID and sessionUUID
func GetUserSession(db *sql.DB, userID string, sessionUUID string)(*UserSession, error){
	dbLogger.Debug("GetUserSession...")
	var err error
	var stmt *sql.Stmt
	var us = new(UserSession)

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return nil, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("SELECT rowid, useruuid, sessionuuid, expiredAt FROM USERSESSION WHERE useruuid = ? and sessionuuid = ?")
	if err != nil {
		dbLogger.Errorf("Failed preparing stmt: %v", err)
		return nil, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	var expiredAtBytes []byte
	if err := stmt.QueryRow(userID, sessionUUID).Scan(&us.RowID, &us.UserUUID, &us.SessionUUID, &expiredAtBytes); err != nil {
		dbLogger.Errorf("Failed getting user session: %v", err)
		return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
	}
	us.ExpiredAt = string(expiredAtBytes)
	dbLogger.Debugf("Get user session: %#v", *us)
	return us, nil
}
