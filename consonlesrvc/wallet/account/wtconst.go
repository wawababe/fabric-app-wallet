package account

import (
	"github.com/op/go-logging"
	util "baas/app-wallet/consonlesrvc/common"
)

var wtLogger *logging.Logger = util.NewLogger("wallet")


type ChaincodeID struct {
	Path string `json:"path,omitempty"`
	Name string `json:"name"`
}
type CtorMsg struct {
	Function string `json:"function"`
	Args []string `json:"args"`
}
type Params struct {
	Type int `json:"type"`
	ChaincodeID ChaincodeID `json:"chaincodeID"`
	CtorMsg CtorMsg `json:"ctorMsg"`
	SecureContext string `json:"securecontext, omitempty"`
}
type PeerReq struct {
	JsonRPC string `json:"jsonrpc"`
	Method string `json:"method"`
	Params Params `json:"params"`
	ID int `json:"id,omitempty"`
}
type InvokeRes struct {
	JsonRPC string `json:"jsonrpc"`
	Result util.BaseResponse `json:"result"`
	ID int `json:"id"`
}
