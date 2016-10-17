package database

import (
	"testing"
	"hash"
	"crypto/sha256"
	"database/sql"
	"io"
	"fmt"
	"github.com/hyperledger/fabric/core/util"
	_ "github.com/go-sql-driver/mysql"
)

func TestAddUser(t *testing.T) {
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		dbLogger.Fatal(err)
	}
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
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		dbLogger.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		dbLogger.Error(ERROR_DB_NOT_CONNECTED)
	}

	var name = "lolshi"
	_, err = GetUserByName(db, name)
	if err != nil {
		t.Errorf("name %s failed: %v", name, err)
	}

}
