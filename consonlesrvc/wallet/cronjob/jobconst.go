package cronjob

import (
	util "baas/app-wallet/consonlesrvc/common"
)

type ChaincodeID struct {
	Path string `json:"path,omitempty"`
	Name string `json:"name"`
}
type CtorMsg struct {
	Function string   `json:"function"`
	Args     []string `json:"args"`
}
type Params struct {
	Type          int         `json:"type"`
	ChaincodeID   ChaincodeID `json:"chaincodeID"`
	CtorMsg       CtorMsg     `json:"ctorMsg"`
	SecureContext string      `json:"securecontext, omitempty"`
}

type PeerInvokeReq struct {
	JsonRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
	ID      int    `json:"id,omitempty"`
}
type PeerInvokeRes struct {
	JsonRPC string            `json:"jsonrpc"`
	Result  util.BaseResponse `json:"result"`
	ID      int               `json:"id"`
}


func NewPeerInvokeReq(funcname string, args []string)(*PeerInvokeReq){
	return &PeerInvokeReq{
		JsonRPC: "2.0",
		Method:  "invoke",
		Params: Params{
			Type: 1,
			ChaincodeID: ChaincodeID{
				Name: "wallet",
			},
			CtorMsg: CtorMsg{
				Function: funcname,
				Args: args,
			},
			SecureContext: "diego",
		},
		ID: 1,
	}
}
