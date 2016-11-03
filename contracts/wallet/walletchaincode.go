package wallet

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"strconv"
	"fmt"
	"encoding/json"
	"time"
)

const (
	DATETIME_FORMAT = "2006-01-02 15:04:05"
)

type WalletChaincode struct {
}


func (w *WalletChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string)([]byte, error){
	return nil, nil
}

func (w *WalletChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string)([]byte, error){
	if function == "createaccount" {
		return w.createAccount(stub, args)
	}else if function == "accounttransfer" {
		return w.accountTransfer(stub, args)
	}
	return nil, fmt.Errorf("unsupported function for invoke: %s", function)
}

func (w *WalletChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string)([]byte, error){
	if function == "getaccount" {
		return w.getAccount(stub, args)
	}else if function == "gettransaction" {
		return w.getTransaction(stub, args)
	}
	return nil, fmt.Errorf("unsupported function for query: %s", function)
}

// args[0]: useruuid
// args[1]: accountuuid
// args[2]: amount
func (w *WalletChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string)([]byte, error){
	var act *Account = new(Account)
	var err error
	if len(args) != 3 {
		wtLogger.Error("incorrect number of arguments. expecting 3")
		return nil, errors.New("incorrect number of arguments. expecting 3")
	}
	act.UserUUID = args[0]
	act.AccountUUID = args[1]
	if act.Amount, err = strconv.ParseInt(args[2], 10, 64); err != nil {
		wtLogger.Errorf("failed to parse amount %s as int64", args[2])
		return nil, fmt.Errorf("failed to parse amount %s as int64", args[2])
	}

	if err = act.Create(stub); err != nil {
		wtLogger.Errorf("failed to create account %#v: %v", *act, err)
		return nil, fmt.Errorf("failed to create account %#v: %v", *act, err)
	}

	wtLogger.Debugf("succeeded in creating account %#v", *act)
	return json.Marshal(act)
}


// args[0]: txuuid
// args[1]: payer accountuuid
// args[2]: payee accountuuid
// args[3]: amount
func (w *WalletChaincode) accountTransfer(stub shim.ChaincodeStubInterface, args []string)([]byte, error){
	var err error
	if len(args) != 4 {
		wtLogger.Debug("incorrect number of arguments for accountTransfer. Expecting 4")
		return nil, errors.New("incorrect number of arguments for accountTransfer. Expecting 4")
	}
	if len(args[0])==0 || len(args[1])==0 || len(args[2])==0 || len(args[3])==0 {
		wtLogger.Errorf("incorrect arguments. expecting not empty")
		return nil, errors.New("incorrect arguments. expecting not empty")
	}

	var txuuid string = args[0]
	var payeruuid string = args[1]
	var payeeuuid string = args[2]
	var payamount int64
	if payamount, err = strconv.ParseInt(args[3], 10, 64); err != nil {
		wtLogger.Errorf("failed to parse transfer amount %s as int64", args[3])
		return nil, fmt.Errorf("failed to parse transfer amount %s as int64", args[3])
	}

	var payer *Account = new(Account)
	var payee *Account = new(Account)
	if !payer.IsAccountExist(stub, payeruuid) {
		wtLogger.Errorf("payer with accountuuid %s not exists", payeruuid)
		return nil, fmt.Errorf("payer with accountuuid %s not exists", payeeuuid)
	}
	if !payee.IsAccountExist(stub, payeeuuid){
		wtLogger.Errorf("payee with accountuuid %s not exists", payeeuuid)
		return nil, fmt.Errorf("payee with accountuuid %s not exists", payeeuuid)
	}


	if payer.Amount < payamount {
		wtLogger.Errorf("payer's account only have %d, not enough to pay %d", payer.Amount, payamount)
		return nil, fmt.Errorf("payer's account only have %d, not enough to pay %d", payer.Amount, payamount)
	}

	payer.Amount -= payamount
	payee.Amount += payamount
	if err = payer.putAccount(stub); err != nil {
		wtLogger.Errorf("failed to modify payer's account: %v", err)
		return nil, fmt.Errorf("failed to modify payer's account: %v", err)
	}

	if err = payee.putAccount(stub); err != nil {
		wtLogger.Errorf("failed to modify payee's account: %v", err)
		payer.Amount += payamount
		payer.putAccount(stub)
		return nil, fmt.Errorf("failed to modify payee's account: %v", err)
	}

	var tx *Transaction = new(Transaction)
	tx.TxUUID = txuuid
	tx.Payeruuid = payeruuid
	tx.Payeeuuid = payeeuuid
	tx.BC_txuuid = stub.GetTxID()
	tx.TxTime = time.Now().Format(DATETIME_FORMAT)
	tx.Status = "fin"

	if err = tx.PutTransaction(stub); err != nil {
		wtLogger.Fatalf("failed to put transaction %#v into ledger: %v", *tx, err)
		payer.Amount += payamount
		payee.Amount -= payamount
		payer.putAccount(stub)
		payee.putAccount(stub)
		return nil, fmt.Errorf("failed to put transaction %#v into ledger: %v", *tx, err)
	}
	wtLogger.Debugf("succeeded in transfer %d money from payer %s to payee %s", payamount, payeruuid, payeeuuid)
	return nil, nil
}

// args[0]: accountuuid
func (w *WalletChaincode) getAccount(stub shim.ChaincodeStubInterface, args []string)([]byte, error){
	if len(args) != 1 || len(args[0])==0 {
		wtLogger.Error("incorrect argument for getAccount. expecting 1, not null")
		return nil, errors.New("incorrect argument for getAccount. expecting 1, not null")
	}

	var account *Account = new(Account)
	var err error
	var accountuuid string = args[0]
	if err = account.GetAccount(stub, accountuuid); err != nil {
		wtLogger.Errorf("failed to get account %s: %v", accountuuid, err)
		return nil, fmt.Errorf("failed to get account %s: %v", accountuuid, err)
	}
	wtLogger.Debug("succeeded in getting account %#v", *account)
	return json.Marshal(account)
}


// args[0]: transactionuuid
func (w *WalletChaincode) getTransaction(stub shim.ChaincodeStubInterface, args []string)([]byte, error){
	if len(args) != 1 || len(args[0])==0 {
		wtLogger.Error("incorrect argument for getTransaction. expecting 1, not null")
		return nil, errors.New("incorrect argument for getTransaction. expecting 1, not null")
	}

	var tx *Transaction = new(Transaction)
	var err error
	var txuuid string = args[0]
	if err = tx.GetTransaction(stub, txuuid); err != nil {
		wtLogger.Errorf("failed to get transaction %s: %v", txuuid, err)
		return nil, fmt.Errorf("failed to get transaction %s: %v", txuuid, err)
	}
	wtLogger.Debugf("succeeded in getting transaction %#v", *tx)
	return json.Marshal(tx)
}
