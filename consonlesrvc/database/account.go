package database

import (
	"database/sql"
	"errors"
	"fmt"
//	"crypto/md5"
)

type Account struct {
	RowID int64 `json:"rowid"`
	AccountUUID string `json:"accountuuid"`
	UserUUID string `json:"useruuid"`
	AccountName string `json:"accountname"`
	AccountID string `json:"accountid"`
	Amount int64 `json:"amount"`
	BC_TXUUID string `json:"bc_txuuid"`
	Status string `json:"status"`
}

func AddAccount(db *sql.DB, u *Account)(int64, error){
	//md5.Sum()
	dbLogger.Debug("AddAccount...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("INSERT INTO account(accountuuid, useruuid, accountname, accountid, amount, bc_txuuid, status) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	addResult, err = stmt.Exec(u.AccountUUID, u.UserUUID, u.AccountName, u.AccountID, u.Amount, u.BC_TXUUID, u.Status)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
}


// GetAccount: get user account by account uuid
func GetAccount(db *sql.DB, accountuuid string) (*Account, error) {
	dbLogger.Debug("GetAccount...")
	var account = new(Account)
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return nil, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("SELECT rowid, accountuuid, useruuid, accountname, accountid, amount, bc_txuuid, status FROM account WHERE accountuuid = ? and deleted = 0")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return nil, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(accountuuid).Scan(&account.RowID, &account.AccountUUID, &account.UserUUID, &account.AccountName, &account.AccountID, &account.Amount, &account.BC_TXUUID, &account.Status); err != nil {
		dbLogger.Errorf("Failed getting account by accountuuid: %v", accountuuid, err)
		return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
	}
	dbLogger.Debugf("Get account by accountuuid %s: \n%#v", accountuuid, *account)

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

	if rows, err = db.Query("SELECT rowid, accountuuid, useruuid, accountname, accountid, amount, bc_txuuid, status FROM account where useruuid = ?", &useruuid); err != nil {
		dbLogger.Errorf("Failed getting accounts by useruuid %s : %v", useruuid, err)
		return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ac *Account = new(Account)
		if err := rows.Scan(&ac.RowID, &ac.AccountUUID, &ac.UserUUID, &ac.AccountName, &ac.AccountID, &ac.Amount, &ac.BC_TXUUID, &ac.Status); err != nil {
			dbLogger.Fatal(err)
			continue //just ignore scan error
		}
		dbLogger.Debugf("useruuid %s has account: %#v", useruuid, *ac)
		accounts = append(accounts, ac)
	}


	return accounts, nil
}

// GetAccount: get user account by accountid
func GetAccountByAccountID(db *sql.DB, accountid string) (*Account, error) {
	dbLogger.Debug("GetAccountByAccountID...")
	var account = new(Account)
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return nil, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("SELECT rowid, accountuuid, useruuid, accountname, accountid, amount, bc_txuuid, status FROM account WHERE accountid = ? and deleted = 0")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return nil, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(accountid).Scan(&account.RowID, &account.AccountUUID, &account.UserUUID, &account.AccountName, &account.AccountID, &account.Amount, &account.BC_TXUUID, &account.Status); err != nil {
		dbLogger.Errorf("Failed getting account by accountuuid: %v", accountid, err)
		return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
	}
	dbLogger.Debugf("Get account by accountid %s: \n%#v", accountid, *account)

	return account, nil
}

/*
func GetAccountByName(db *sql.DB, name string)(*Account, error){
	dbLogger.Debug("GetAccountByName...")
	var account = new(Account)
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return nil, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("SELECT rowid, accountuuid, useruuid, amount, accountname, bc_txuuid, status FROM account WHERE accountname = ? and deleted = 0")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return nil, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(name).Scan(&account.RowID, &account.AccountUUID, &account.UserUUID, &account.Amount, &account.AccountName, &account.BC_TXUUID, &account.Status); err != nil {
		dbLogger.Errorf("Failed getting account by name %s", name, err)
		return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
	}
	dbLogger.Debugf("Get account by name %s: \n%#v", name, *account)

	return account, nil
}
*/

func DeleteAccount(db *sql.DB, u *Account)(int64, error){
	dbLogger.Debug("DeleteAccount...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("UPDATE account SET deleted = 1 and status = ? WHERE accountuuid = ?")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	addResult, err = stmt.Exec(u.Status, u.AccountUUID)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
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

	stmt, err = db.Prepare("UPDATE account SET amount = ?, bc_txuuid = ?, status = ? WHERE accountuuid = ?")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	addResult, err = stmt.Exec(u.Amount, u.BC_TXUUID, u.Status, u.AccountUUID)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
}