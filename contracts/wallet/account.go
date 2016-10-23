package wallet

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)


// todo: whether it's a need to add additional information, such as bankname
type Account struct {
	AccountUUID string `json:"accountuuid"`
	UserUUID    string `json:"useruuid"`
	Amount int64 `json:"amount"`
}

// Create: create a new account
func (t *Account) Create(stub shim.ChaincodeStubInterface) (error) {
	wtLogger.Debug("start to create account....")
	defer func(){
		wtLogger.Debug("start to create account...Done!")
	}()

	var err error

	if len(t.UserUUID) == 0 || len(t.AccountUUID) == 0 {
		wtLogger.Errorf("account %#v not valid", *t)
		return fmt.Errorf("account %#v not valid", *t)
	}

	var dupAcc *Account = new(Account)
	if err = dupAcc.GetAccount(stub, t.AccountUUID); err == nil {
		wtLogger.Errorf("illegal to create duplicate account %#v", *dupAcc)
		return fmt.Errorf("illegal to create duplicate account %#v", *dupAcc)
	}

	if err = t.putAccount(stub); err != nil {
		wtLogger.Errorf("failed putaccount: %v", err)
		return fmt.Errorf("failed putaccount: %v", err)
	}
	wtLogger.Debugf("successed in creating account %#v", *t)
	return nil
}

// todo: should not delete account out of security
func (t *Account) Delete(stub shim.ChaincodeStubInterface, account string) error {
	return nil
}



func (t *Account) putAccount(stub shim.ChaincodeStubInterface) error {
	acBytes, err := json.Marshal(*t)
	if err != nil {
		wtLogger.Fatalf("failed marshalling account %#v as bytes: %v", *t, err)
		return fmt.Errorf("failed marshalling account %#v as bytes: %v", *t, err)
	}
	wtLogger.Debugf("marshalled account %#v as bytes", *t)

	if err := stub.PutState(t.buildKey(t.AccountUUID), acBytes); err != nil {
		wtLogger.Fatalf("failed putting account %#v into ledger: %v", *t, err)
		return fmt.Errorf("failed putting account %#v into ledger: %v", *t, err)
	}
	wtLogger.Debugf("successed in putting account %#v with key %s into ledger", *t, t.buildKey(t.AccountUUID))

	return nil
}

func (t *Account) buildKey(accountuuid string) (key string) {
	key = "acn-" + accountuuid
	return
}

func (t *Account) GetAccount(stub shim.ChaincodeStubInterface, accountuuid string) (error) {
	wtLogger.Debug("start to GetAccount...")
	defer func(){
		wtLogger.Debug("start to GetAccount...Done!")
	}()
	var acBytes []byte
	var err error
	if acBytes, err = stub.GetState(t.buildKey(accountuuid)); err != nil {
		wtLogger.Warningf("failed getting state for account %s: %v", accountuuid, err)
		return fmt.Errorf("failed getting state for account %s: %v", accountuuid, err)
	}

	if err = json.Unmarshal(acBytes, t); err != nil {
		wtLogger.Warningf("failed unmarshalling acbytes into account: %v", err)
		return fmt.Errorf("failed unmarshalling acbytes into account: %v", err)
	}
	wtLogger.Debugf("got account %#v", *t)
	return nil
}

func (t *Account) IsAccountExist(stub shim.ChaincodeStubInterface, accountuuid string)(bool){
	if err := t.GetAccount(stub, accountuuid); err != nil {
		return false // account not exist
	}
	wtLogger.Debugf("account with uuid %s exists: %#v", accountuuid, *t)
	return true
}
