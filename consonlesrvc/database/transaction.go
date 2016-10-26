package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type Transaction struct {
	RowID int64 `json:"rowid"`
	TxUUID string `json:"txuuid"`
	PayerUUID string `json:"payeruuid"` //payer's useruuid
	PayeeUUID string `json:"payeeuuid"` //payee's useruuid
	PayerAccountID string `json:"payeraccountid"`//payer accountid
	PayeeAccountID string `json:"payeeaccountid"`//payee accountid
	Amount int64 `json:"amount"`
	BC_txuuid string `json:"bc_txuuid"`
	BC_blocknum int64 `json:"bc_blocknum"`
	Status string `json:"status"`//pending, fin, failed
}

func AddTransaction(db *sql.DB, t *Transaction)(int64, error){
	dbLogger.Debug("AddTransaction...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("INSERT INTO transaction(txuuid, payeruuid, payeeuuid, payeraccountid, payeeaccountid, amount, bc_txuuid, bc_blocknum, status) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	addResult, err = stmt.Exec(t.TxUUID, t.PayerUUID, t.PayeeUUID, t.PayerAccountID, t.PayeeAccountID, t.Amount, t.BC_txuuid, t.BC_blocknum, t.Status)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
}


// GetTransaction: get transaction by txuuid
func GetTransaction(db *sql.DB, txuuid string)(*Transaction, error){
	dbLogger.Debug("GetTransaction...")
	var tx = new(Transaction)
	var err error
	var stmt *sql.Stmt

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return nil, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("SELECT rowid, txuuid, payeruuid, payeeuuid, payeraccountid, payeeaccountid, amount, bc_txuuid, bc_blocknum, status FROM transaction WHERE txuuid = ? and deleted = 0")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return nil, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(txuuid).Scan(&tx.RowID, &tx.TxUUID, &tx.PayerUUID, &tx.PayeeUUID, &tx.PayerAccountID, &tx.PayeeAccountID, &tx.Amount, &tx.BC_txuuid, &tx.BC_blocknum, &tx.Status); err != nil {
		dbLogger.Errorf("Failed getting transaction by txuuid %s: %v", txuuid, err)
		return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
	}
	dbLogger.Debugf("Get transaction by txuuid %s: \n%+v", txuuid, *tx)

	return tx, nil
}


func GetTransactionsByPayeruuid(db *sql.DB, payeruuid string)([]*Transaction, error){
	dbLogger.Debug("GetTransactionsByPayeruuid...")
	var err error
	var rows *sql.Rows
	var txs []*Transaction

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return nil, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	if rows, err = db.Query("SELECT rowid, txuuid, payeruuid, payeeuuid, payeraccountid, payeeaccountid, amount, bc_txuuid, bc_blocknum, status FROM transaction where payeruuid = ? and deleted = 0", &payeruuid); err != nil {
		dbLogger.Errorf("Failed getting transactions by payeruuid %s : %v", payeruuid, err)
		return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
	}
	if err != nil {
		dbLogger.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var t *Transaction = new(Transaction)
		if err := rows.Scan(&t.RowID, &t.TxUUID, &t.PayerUUID, &t.PayeeUUID, &t.PayerAccountID, &t.PayeeAccountID, &t.Amount, &t.BC_txuuid, &t.BC_blocknum, &t.Status); err != nil {
			dbLogger.Fatal(err)
		}
		dbLogger.Debugf("payer %s has transaction: %#v", payeruuid, *t)
		txs = append(txs, t)
	}


	return txs, nil

}

func GetTransactionsByPayeeuuid(db *sql.DB, payeeuuid string)([]*Transaction, error){
	dbLogger.Debug("GetTransactionsByPayeeuuid...")
	var err error
	var rows *sql.Rows
	var txs []*Transaction

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return nil, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	if rows, err = db.Query("SELECT rowid, txuuid, payeruuid, payeeuuid, payeraccountid, payeeaccountid, amount, bc_txuuid, bc_blocknum, status FROM transaction where payeeuuid = ? and deleted = 0", &payeeuuid); err != nil {
		dbLogger.Errorf("Failed getting transactions by payeeuuid %s : %v", payeeuuid, err)
		return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
	}
	if err != nil {
		dbLogger.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var t *Transaction = new(Transaction)
		if err := rows.Scan(&t.RowID, &t.TxUUID, &t.PayerUUID, &t.PayeeUUID, &t.PayerAccountID, &t.PayeeAccountID, &t.Amount, &t.BC_txuuid, &t.BC_blocknum, &t.Status); err != nil {
			dbLogger.Fatal(err)
		}
		dbLogger.Debugf("payee %s has transaction: %#v", payeeuuid, *t)
		txs = append(txs, t)
	}


	return txs, nil

}

func UpdateTransaction(db *sql.DB, t *Transaction)(int64, error){
	dbLogger.Debug("UpdateTransaction...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("UPDATE transaction SET bc_txuuid = ?, bc_blocknum = ?, status = ? WHERE txuuid = ?")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	addResult, err = stmt.Exec(t.BC_txuuid, t.BC_blocknum, t.Status, t.TxUUID)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
}

func DeleteTransaction(db *sql.DB, t *Transaction)(int64, error){
	dbLogger.Debug("UpdateTransaction...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("UPDATE transaction SET deleted = 1 and status = ? WHERE txuuid = ?")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	addResult, err = stmt.Exec(t.Status, t.TxUUID)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
}