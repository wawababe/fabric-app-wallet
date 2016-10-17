package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type Account struct {
	RowID int64
	AccountUUID string
	UserUUID string
	Amount int64
}

func AddAccount(db *sql.DB, u *Account)(int64, error){
	dbLogger.Debug("AddAccount...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("INSERT INTO account(accountuuid, useruuid, amount) VALUES(?, ?, ?)")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	addResult, err = stmt.Exec(u.AccountUUID, u.UserUUID, u.Amount)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
}


// GetAccount: get user account by useruuid and account uuid
func GetAccount(db *sql.DB, useruuid string, accountuuid string) (*Account, error) {
	dbLogger.Debug("GetAccount...")
	var account = new(Account)
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return nil, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("SELECT rowid, accountuuid, useruuid, amount FROM account WHERE useruuid = ? and accountuuid = ? and deleted = 0")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return nil, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(useruuid, accountuuid).Scan(&account.RowID, &account.AccountUUID, &account.UserUUID, &account.Amount); err != nil {
		dbLogger.Errorf("Failed getting account by useruuid %s and accountuuid: %v", useruuid, accountuuid, err)
		return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
	}
	dbLogger.Debugf("Get account by useruuid %s and accountuuid %s: \n%#v", useruuid, accountuuid, *account)

	return account, nil
}

func GetAccountsByUseruuid(db *sql.DB, useruuid string)([]*Account, error){
	dbLogger.Debug("GetAccountsByUseruuid...")
	var err error
	var rows *sql.Rows
	var accounts []*Account

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return nil, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	if rows, err = db.Query("SELECT rowid, accountuuid, useruuid, amount FROM account where useruuid = ?", &useruuid); err != nil {
		dbLogger.Errorf("Failed getting accounts by useruuid %s : %v", useruuid, err)
		return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
	}
	if err != nil {
		dbLogger.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var ac *Account = new(Account)
		if err := rows.Scan(&ac.RowID, &ac.AccountUUID, &ac.UserUUID, &ac.Amount); err != nil {
			dbLogger.Fatal(err)
		}
		dbLogger.Debugf("useruuid %s has account: %#v", useruuid, *ac)
		accounts = append(accounts, ac)
	}


	return accounts, nil
}


func UpdateAccount(db *sql.DB, u *Account)(int64, error){
	dbLogger.Debug("UpdateAccount...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("UPDATE account SET amount = ? WHERE accountuuid = ?")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	addResult, err = stmt.Exec(u.Amount, u.AccountUUID)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
}