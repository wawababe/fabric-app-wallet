package database

import (
	"database/sql"
	"errors"
	"fmt"
)


type User struct {
	RowID    int64
	UserUUID string
	Password string
	Username string
}

func AddUser(db *sql.DB, u *User) (int64, error) {
	dbLogger.Debug("AddUser...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("INSERT INTO USER(useruuid, password, username) VALUES(?, ?, ?)")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	addResult, err = stmt.Exec(u.UserUUID, u.Password, u.Username)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err

}

func GetUserByName(db *sql.DB, name string) (*User, error) {
	dbLogger.Debug("GetUserByName...")
	var user = new(User)
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return nil, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("SELECT rowid, useruuid, password, username FROM user WHERE username = ? and deleted = 0")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return nil, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(name).Scan(&user.RowID, &user.UserUUID, &user.Password, &user.Username); err != nil {
		dbLogger.Errorf("Failed getting user %s: %v", name, err)
		return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
	}
	dbLogger.Debugf("Get user by name %s: \n%+v", name, *user)

	return user, nil
}
