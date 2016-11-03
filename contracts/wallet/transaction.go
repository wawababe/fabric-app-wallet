package wallet

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type Transaction struct {
	TxUUID string  `json:"txuuid"`
	Payeruuid  string `json:"payeruuid"`
	Payeeuuid string `json:"payeeuuid"`
	Amount string  `json:"amount"`
	TxTime string  `json:"txtime"`
	BC_txuuid string `json:"bc_txuuid"`
	Status string `json:"status"`
}


// write the transaction into blockchain
func (t *Transaction) PutTransaction(stub shim.ChaincodeStubInterface) (error) {
	wtLogger.Debug("start to PutTransaction...")
	defer func(){
		wtLogger.Debug("start to PutTransaction...Done!")
	}()
	txBytes, err := json.Marshal(*t)
	if err != nil {
		wtLogger.Fatalf("failed marshalling transaction %#v as bytes: %v", *t, err)
		return fmt.Errorf("failed marshalling transaction %#v: %v", *t, err)
	}
	wtLogger.Debugf("marshalled transaction %#v as bytes", *t)

	if err := stub.PutState(t.buildKey(t.TxUUID), txBytes); err != nil {
		wtLogger.Fatalf("failed putting transaction %#v into ledger: %v", *t, err)
		return fmt.Errorf("failed putting transaction %#v into ledger: %v", *t, err)
	}
	wtLogger.Debugf("succeeded in putting transaction %#v with key %s into ledger", *t, t.buildKey(t.TxUUID))

	return nil
}

func (t *Transaction) buildKey(txUUID string) (key string) {
	key = "tx-" + txUUID
	return
}

func (t *Transaction) GetTransaction(stub shim.ChaincodeStubInterface, txuuid string) error {
	wtLogger.Debug("start to GetTransaction...")
	defer func(){
		wtLogger.Debug("start to GetTransaction...Done!")
	}()
	var txBytes []byte
	var err error
	if txBytes, err = stub.GetState(t.buildKey(txuuid)); err != nil {
		wtLogger.Warningf("failed getting state for transaction %s: %v", txuuid, err)
		return fmt.Errorf("failed getting state for transaction %s: %v", txuuid, err)
	}

	if err = json.Unmarshal(txBytes, t); err != nil {
		wtLogger.Warningf("failed unmarshalling txbytes into transaction: %v", err)
		return fmt.Errorf("failed unmarshalling txbytes into transadtion: %v", err)
	}
	wtLogger.Debugf("got transaction %#v", *t)
	return nil
}

//
