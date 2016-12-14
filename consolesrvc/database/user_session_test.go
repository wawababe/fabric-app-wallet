package database

import (
	"testing"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	util "baas/app-wallet/consolesrvc/common"
)


func TestAddUserSession(t *testing.T) {
	var db *sql.DB = GetDB()
	defer db.Close()

	var testNum = 2
	var usersessions = make([]UserSession, testNum)
	for i, _ := range usersessions {
		usersessions[i].SessionUUID = util.GenerateUUID()
		usersessions[i].UserUUID = util.GenerateUUID()
		usersessions[i].AddExpiredTimeByDays(1)
		dbLogger.Debugf("user session: %#v", usersessions[i])
	}
	usersessions[0].UserUUID = "5cdb617c-2712-480a-a02b-facd8c86e579"
	usersessions[1].SessionUUID = usersessions[0].SessionUUID

	var tests = []struct{
		newline bool
		sep string
		arg *UserSession
		want int64
	}{
		{false, " ", &usersessions[0], 1},
		{true, " ", &usersessions[1], 0},
	}

	for i, testitem := range tests {
		rowsAff, _ := AddUserSession(db, testitem.arg)
		if rowsAff != testitem.want {
			t.Errorf("Test #%d: Add session item %#v, affected rows = %d, but want %d", i, testitem.arg, rowsAff, testitem.want)
		}
	}

}

func TestGetUserSession(t *testing.T) {
	var db *sql.DB = GetDB()
	defer db.Close()
	var err error
	var us *UserSession

	var useruuid, sessionuuid string
	useruuid = "5cdb617c-2712-480a-a02b-facd8c86e579"
	sessionuuid = "59336beb-9b1e-4467-8e9a-c88dd553484f"
	us, err = GetUserSession(db, useruuid, sessionuuid)
	if us == nil || err != nil{
		t.Errorf("Failed retrieving user session")
	}
	dbLogger.Debugf("Get user session: %#v", *us)
}
