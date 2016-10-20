package main

import (
	"baas/app-wallet/contracts/wallet"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/hyperledger/fabric/core/crypto/primitives"
)

func main() {
	//primitives.SetSecurityLevel("SHA3", 256)
	err := shim.Start(new(wallet.WalletChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err);
	}

}