package database

import (
	"testing"
	"database/sql"
	"github.com/hyperledger/fabric/core/util"
)

func TestAddAccount(t *testing.T) {
	var db *sql.DB
	var err error
	if db, err = sql.Open("mysql", DSN); err != nil {
		dbLogger.Fatal("Failed opening database")
	}

	var testNum = 2

	var accounts = make([]Account, testNum)
	for i, _ := range accounts {
		accounts[i].AccountUUID = util.GenerateUUID()
		accounts[i].UserUUID = util.GenerateUUID()
		accounts[i].Amount = 200
		dbLogger.Debugf("user account: %#v", accounts[i])
	}
	accounts[0].UserUUID = "5cdb617c-2712-480a-a02b-facd8c86e579"
	accounts[1].AccountUUID = accounts[0].AccountUUID

	var tests = []struct{
		newline bool
		sep string
		arg *Account
		want int64
	}{
		{false, " ", &accounts[0], 1},
		{true, " ", &accounts[1], 0},
	}

	for i, testitem := range tests {
		rowsAff, _ := AddAccount(db, testitem.arg)
		if rowsAff != testitem.want {
			t.Errorf("Test #%d: Add account %#v, affected rows = %d, but want %d", i, testitem.arg, rowsAff, testitem.want)
		}
	}
}


func TestGetAccountByUserUUID(t *testing.T) {
	var db *sql.DB
	var err error
	var us *Account
	if db, err = sql.Open("mysql", DSN); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
	}

	var useruuid, accountuuid string
	useruuid = "5cdb617c-2712-480a-a02b-facd8c86e579"
	accountuuid = "b7e97e66-dba8-4cf7-af2f-fe17ee7e7c03"
	us, err = GetAccount(db, useruuid, accountuuid)
	if us == nil || err != nil{
		t.Error("Failed retrieving user account: %v", err)
	}
	dbLogger.Debugf("Get user account: %#v", *us)
}


func TestGetAccountsByUseruuid(t *testing.T) {
	var db *sql.DB
	var err error
	var accounts []*Account
	if db, err = sql.Open("mysql", DSN); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
	}

	var useruuid string = "5cdb617c-2712-480a-a02b-facd8c86e579"
	accounts, err = GetAccountsByUseruuid(db, useruuid)
	if err != nil {
		t.Errorf("Failed retrieving user accounts by useruuid %s: %v", useruuid, err)
	}
	for i, account := range accounts {
		dbLogger.Debugf("Accounts #%d: %v", i, *account)
	}

}

func TestUpdateAccount(t *testing.T) {
	var db *sql.DB
	var err error
	var us *Account
	if db, err = sql.Open("mysql", DSN); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
	}

	var useruuid, accountuuid string
	useruuid = "5cdb617c-2712-480a-a02b-facd8c86e579"
	accountuuid = "b7e97e66-dba8-4cf7-af2f-fe17ee7e7c03"
	us, err = GetAccount(db, useruuid, accountuuid)
	if us == nil || err != nil{
		t.Error("Failed retrieving user account: %v", err)
	}
	dbLogger.Debugf("Get user account: %#v", *us)

	us.Amount -= 10

	var affectedrows int64 = 0
	affectedrows, err = UpdateAccount(db, us)
	if affectedrows != 1 {
		t.Errorf("Failed updating account to %v\n err: %v", *us, err)
	}

	dbLogger.Debugf("Updated user account: %#v", *us)


}