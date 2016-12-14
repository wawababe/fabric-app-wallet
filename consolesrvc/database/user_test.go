package database

import (
	"testing"
	"hash"
	"crypto/sha256"
	"io"
	"fmt"
	util "baas/app-wallet/consolesrvc/common"
	_ "github.com/go-sql-driver/mysql"
)

func TestAddUser(t *testing.T) {
	db := GetDB()
	defer db.Close()

	if err := db.Ping(); err != nil {
		dbLogger.Error(ERROR_DB_NOT_CONNECTED)
	}

	var csum hash.Hash
	csum = sha256.New()
	io.WriteString(csum, "123")
	s := fmt.Sprintf("%x", csum.Sum(nil))

	dbLogger.Debug("Hashed password: " + s)

	var user = &User{
		UserUUID: util.GenerateUUID(),
		Username: "lolshi",
		Password: s,
	}
	if rowsAff, err := AddUser(db, user); err != nil {
		t.Errorf("Failed adding user: rowsaffected = %d; %v", rowsAff, err)
	}

}

func TestGetUserByName(t *testing.T) {
	db := GetDB()
	defer db.Close()
	var err error

	if err = db.Ping(); err != nil {
		dbLogger.Error(ERROR_DB_NOT_CONNECTED)
	}

	var name = "lolshi"
	_, err = GetUserByName(db, name)
	if err != nil {
		t.Errorf("name %s failed: %v", name, err)
	}

}
