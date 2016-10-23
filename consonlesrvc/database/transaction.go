package database

import (
	"database/sql"
	"errors"
	"fmt"
)

type Transaction struct {
	RowID int64
	TXUUID string
	PayerUUID string //payer account uuid
	PayeeUUID string //payee account uuid
	Amount int64
	BC_txuuid string
	BC_blocknum int64
	Status string //pending, fin, failed
}

func AddTransaction(db *sql.DB, u *Transaction)(int64, error){
	dbLogger.Debug("AddTransaction...")
	var err error
	var stmt *sql.Stmt
	var addResult sql.Result
	var affectedRows int64

	if err := db.Ping(); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
		return 0, errors.New(ERROR_DB_NOT_CONNECTED)
	}

	stmt, err = db.Prepare("INSERT INTO transaction(txuuid, payeruuid, payeeuuid, amount, bc_txuuid, bc_blocknum, status) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return 0, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	addResult, err = stmt.Exec(u.TXUUID, u.PayerUUID, u.PayeeUUID, u.Amount, u.BC_txuuid, u.BC_blocknum, u.Status)
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

	stmt, err = db.Prepare("SELECT rowid, txuuid, payeruuid, payeeuuid, amount, bc_txuuid, bc_blocknum, status FROM transaction WHERE txuuid = ?")
	if err != nil {
		dbLogger.Errorf("Failed preparing statement: %v", err)
		return nil, fmt.Errorf(ERROR_DB_PREPARED + ": %v", err)
	}
	defer stmt.Close()

	if err := stmt.QueryRow(txuuid).Scan(&tx.RowID, &tx.TXUUID, &tx.PayerUUID, &tx.PayeeUUID, &tx.Amount, &tx.BC_txuuid, &tx.BC_blocknum, &tx.Status); err != nil {
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

	if rows, err = db.Query("SELECT transaction.rowid, txuuid, payeruuid, payeeuuid, transaction.amount, transaction.bc_txuuid, bc_blocknum, transaction.status FROM transaction INNER JOIN account ON transaction.payeruuid = account.accountuuid where useruuid = ?", &payeruuid); err != nil {
		dbLogger.Errorf("Failed getting transactions by payeruuid %s : %v", payeruuid, err)
		return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
	}
	if err != nil {
		dbLogger.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var t *Transaction = new(Transaction)
		if err := rows.Scan(&t.RowID, &t.TXUUID, &t.PayerUUID, &t.PayeeUUID, &t.Amount, &t.BC_txuuid, &t.BC_blocknum, &t.Status); err != nil {
			dbLogger.Fatal(err)
		}
		dbLogger.Debugf("useruuid %s has transaction: %#v", payeruuid, *t)
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

	if rows, err = db.Query("SELECT transaction.rowid, txuuid, payeruuid, payeeuuid, transaction.amount, transaction.bc_txuuid, bc_blocknum, transaction.status FROM transaction INNER JOIN account ON transaction.payeeuuid = account.accountuuid where useruuid = ?", &payeeuuid); err != nil {
		dbLogger.Errorf("Failed getting transactions by payeeuuid %s : %v", payeeuuid, err)
		return nil, fmt.Errorf(ERROR_DB_QUERY + ": %v", err)
	}
	if err != nil {
		dbLogger.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var t *Transaction = new(Transaction)
		if err := rows.Scan(&t.RowID, &t.TXUUID, &t.PayerUUID, &t.PayeeUUID, &t.Amount, &t.BC_txuuid, &t.BC_blocknum, &t.Status); err != nil {
			dbLogger.Fatal(err)
		}
		dbLogger.Debugf("useruuid %s has transaction: %#v", payeeuuid, *t)
		txs = append(txs, t)
	}


	return txs, nil

}

func UpdateTransaction(db *sql.DB, u *Transaction)(int64, error){
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

	addResult, err = stmt.Exec(u.BC_txuuid, u.BC_blocknum, u.Status, u.TXUUID)
	if err != nil {
		dbLogger.Errorf("Failed executing statement:  %v", err)
		return 0, fmt.Errorf(ERROR_DB_EXECUTE + ": %v", err)
	}

	affectedRows, err = addResult.RowsAffected()
	dbLogger.Debugf("Affected rows: %d", affectedRows)
	return affectedRows, err
}