package blockchain

import (
	"baas/app-wallet/consonlesrvc/auth"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"net/http"
	"fmt"
	"os"
	"strings"
	util "baas/app-wallet/consonlesrvc/common"
)

type TransactionDetailRequest struct {
	authsrvc.AuthRequest
	BC_TXUUID string `json:"bc_txuuid"`
}

type TransactionDetailResponse struct {
	authsrvc.AuthResponse
	TxDetail *BCTransactionMsg `json:"txdetail,omitempty"`
}

type TransactionDetail struct {
}


func (t *TransactionDetail) post(req *TransactionDetailRequest)(*TransactionDetailResponse){
	var res *TransactionDetailResponse = new(TransactionDetailResponse)
	var err error

	if !req.IsRequestValid(&res.AuthResponse) {
		bcLogger.Warningf("request not valid: %#v", *req)
		res.Status = "error"
		res.Message = util.ERROR_UNAUTHORIZED
		res.UserUUID = ""
		return res
	}
	if len(req.BC_TXUUID) == 0 {
		res.Status = "error"
		res.Message = util.ERROR_BADREQUEST + fmt.Sprint(": blockchain transaction uuid should not be empty")
		res.UserUUID = ""
		return res
	}

	var peerAddr string
	if peerAddr = os.Getenv("PEER_ADDRESS"); len(peerAddr) == 0 {
		bcLogger.Fatal("failed getting environmental variable PEER_ADDRESS")
		res.Status = "error"
		res.Message = "failed getting peer address"
		return res
	}
	peerAddr += "/transactions/" + req.BC_TXUUID

	var peerResp *http.Response = new(http.Response)
	if peerResp, err = http.Get(peerAddr); err != nil {
		bcLogger.Errorf("failed querying blockchain transaction %s: %v",  req.BC_TXUUID, err)
		res.Status = "error"
		res.Message = "failed posint createAccount request to peer"
		return res
	}
	defer peerResp.Body.Close()

	if peerResp.StatusCode == http.StatusNotFound {
		bcLogger.Errorf("transaction %s not exist in the blockchain", req.BC_TXUUID)
		res.Status = "error"
		res.Message = util.ERROR_NOTFOUND + fmt.Sprintf(": transaction %s not exist in the blockchain", req.BC_TXUUID)
		return res
	}

	res.TxDetail = new(BCTransactionMsg)
	if err = json.NewDecoder(peerResp.Body).Decode(res.TxDetail); err != nil {
		bcLogger.Error("failed to decode transaction message from blockchain")
		res.Status = "error"
		res.Message = "failed to decode transaction message from blockchain"
	}
	bcLogger.Debugf("decoded transaction message : %#v", *res.TxDetail)
	//var buf = new(bytes.Buffer)
	//io.Copy(buf, peerResp.Body)
	//res.TxDetail = buf.String()
	//res.TxDetail = fmt.Sprintf("%s", buf.String())

	bcLogger.Debugf("got blockchain transaction %s: %#v", req.BC_TXUUID, res.TxDetail)
	return res
}

func TransactionDetailPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params){
	var err error
	var req *TransactionDetailRequest = new(TransactionDetailRequest)
	var res *TransactionDetailResponse
	var resBytes []byte

	w.Header().Set("Content-Type", "application/json")

	if err = r.ParseForm(); err != nil {
		bcLogger.Fatalf("failed to parse request for url %s: %v", r.URL.Path, err)
	}

	req.Username = r.PostForm.Get("username")
	req.SessionID = r.PostForm.Get("sessionid")
	req.AuthToken = r.PostForm.Get("authtoken")
	req.BC_TXUUID = r.PostForm.Get("bc_txuuid")
	bcLogger.Debugf("parsed request for url %s: %#v", r.URL.Path, req)

	var t TransactionDetail
	res = t.post(req)

	if strings.Contains(res.Message, util.ERROR_UNAUTHORIZED){
		w.WriteHeader(http.StatusUnauthorized)
	}else if strings.Contains(res.Message, util.ERROR_BADREQUEST){
		w.WriteHeader(http.StatusBadRequest)
	}else if strings.Contains(res.Message, util.ERROR_NOTFOUND){
		w.WriteHeader(http.StatusNotFound)
	}else if res.Status == "error"{
		w.WriteHeader(http.StatusNotFound)
	}


	resBytes, err = json.Marshal(*res)
	if err != nil {
		bcLogger.Fatalf("failed to marshal response as []byte: %v", err)
	}
	fmt.Fprintf(w, "%s", string(resBytes))
}
