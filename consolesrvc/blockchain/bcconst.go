package blockchain

import (
	"github.com/op/go-logging"
	"baas/app-wallet/consolesrvc/common"
	"github.com/golang/protobuf/ptypes/timestamp"
)

var bcLogger *logging.Logger = common.NewLogger("blockchain")

type BCTransactionMsg struct {
	Type int `json:"type,omitempty"`
	ChaincodeID string `json:"chaincodeID,omitempty"`
	Payload string `json:"payload,omitempty"`
	TxID string `json:"txid,omitempty"`
	Timestamp timestamp.Timestamp `json:"timestamp,omitempty"`
	Nonce string `json:"nonce,omitempty"`
	Cert string `json:"cert,omitempty"`
	Signature string `json:"signature,omitempty"`
}
