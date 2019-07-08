package helpers

import (
	"encoding/hex"
	"encoding/json"
	"math"
	"net"
	"strconv"
	s "strings"
	"time"

	secp "github.com/vncsphere-foundation/secp256k1-go"
)

func PubkeyFromSeckey(privateKey []byte) []byte {
	//

	return secp.PubkeyFromSeckey([]byte(privateKey))
}

func GetRawTransactionType(rawData string) string {
	//

	var tranType TranType
	err := json.Unmarshal([]byte(rawData), &tranType)
	if err != nil {
		//

		return ""

	} else {
		//

		return tranType.TT
	}
}

func ParseHelloTransaction(rawData string) (HelloTransaction, error) {
	//

	var helloTransaction HelloTransaction
	err := json.Unmarshal([]byte(rawData), &helloTransaction)
	return helloTransaction, err
}

func VerifyHelloTransaction(tran HelloTransaction) (ok bool) {
	//

	ok = false

	if len(tran.SENDER) != 66 {
		//

		return
	}

	if net.ParseIP(tran.ADDRESS) == nil {
		//

		return
	}

	if len(tran.TST) != 10 {
		//

		return
	}

	transactionTime, err := strconv.ParseInt(tran.TST, 10, 64)
	if err != nil {
		//

		return
	}

	timestamp := time.Unix(transactionTime, 0)
	if int64(math.Abs(float64(time.Since(timestamp)/time.Second))) > 10 {
		//

		return
	}

	if len(tran.SIGNATURE) != 130 {
		//

		return
	}

	transcationForVerify := HelloTransactionForVerify{tran.TT, tran.SENDER, tran.ADDRESS, tran.TST}
	js, err := json.Marshal(transcationForVerify)
	if err != nil {
		//

		return
	}

	decodedSignature, err := hex.DecodeString(tran.SIGNATURE)
	if err != nil {
		//

		return
	}

	publicKey, err := hex.DecodeString(tran.SENDER)
	if err != nil {
		//

		return
	}

	if secp.VerifySignature(js, decodedSignature, publicKey) != 1 {
		//

		return
	}

	ok = true

	return
}

func ParseSimpleTransaction(rawData string) (SimpleTransaction, error) {
	//

	var simpleTransaction SimpleTransaction
	err := json.Unmarshal([]byte(rawData), &simpleTransaction)
	return simpleTransaction, err
}

func VerifySimpleTransaction(tran SimpleTransaction) (transactionTime int64, statusCode int, ok bool) {
	//

	transactionTime = 0
	ok = false
	statusCode = StatusOk

	if len(tran.SENDER) != 66 {
		//

		statusCode = StatusWrongAttr_SENDER
		return
	}

	if len(tran.RECEIVER) != 66 {
		//

		statusCode = StatusWrongAttr_RECEIVER
		return
	}

	if tran.SENDER == tran.RECEIVER {
		//

		statusCode = StatusDontSendYourself
		return
	}

	if len(tran.TST) != 10 {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	transactionTime, err := strconv.ParseInt(tran.TST, 10, 64)
	if err != nil {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	timestamp := time.Unix(transactionTime, 0)
	if int64(math.Abs(float64(time.Since(timestamp)/time.Second))) > 10 {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	if len(tran.TTOKEN) == 0 {
		//

		statusCode = StatusWrongAttr_TTOKEN
		return
	}

	if len(tran.CTOKEN) == 0 {
		//

		statusCode = StatusWrongAttr_CTOKEN
		return

	} else {
		//

		_, err = strconv.ParseFloat(tran.CTOKEN, 64)
		if err != nil {
			//

			statusCode = StatusWrongAttr_CTOKEN
			return
		}

		if len(tran.CTOKEN[s.Index(tran.CTOKEN, ".")+1:]) > 8 {
			//

			statusCode = StatusWrongAttr_CTOKEN
			return
		}
	}

	if len(tran.SIGNATURE) != 130 {
		//

		statusCode = StatusWrongAttr_Signature
		return
	}

	transcationForVerify := SimpleTransactionForVerify{tran.TT, tran.SENDER, tran.RECEIVER, tran.TTOKEN, tran.CTOKEN, tran.TST}
	js, err := json.Marshal(transcationForVerify)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	decodedSignature, err := hex.DecodeString(tran.SIGNATURE)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	publicKey, err := hex.DecodeString(tran.SENDER)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	if secp.VerifySignature(js, decodedSignature, publicKey) != 1 {
		//

		statusCode = StatusSignVerifyError
		return
	}

	ok = true

	return
}

func ParseApplicantTransaction(rawData string) (ApplicantTransaction, error) {
	//

	var applicantTransaction ApplicantTransaction
	err := json.Unmarshal([]byte(rawData), &applicantTransaction)
	return applicantTransaction, err
}

func VerifyApplicantTransaction(tran ApplicantTransaction) (transactionTime int64, statusCode int, ok bool) {
	//

	transactionTime = 0
	ok = false
	statusCode = StatusOk

	if len(tran.SENDER) != 66 {
		//

		statusCode = StatusWrongAttr_SENDER
		return
	}

	if ip := net.ParseIP(tran.IPADDR); ip == nil {
		//

		statusCode = StatusWrongAttr_IPADDR
		return
	}

	if len(tran.TST) != 10 {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	transactionTime, err := strconv.ParseInt(tran.TST, 10, 64)
	if err != nil {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	timestamp := time.Unix(transactionTime, 0)
	if int64(math.Abs(float64(time.Since(timestamp)/time.Second))) > 10 {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	if len(tran.SIGNATURE) != 130 {
		//

		statusCode = StatusWrongAttr_Signature
		return
	}

	transcationForVerify := ApplicantTransactionForVerify{tran.TT, tran.SENDER, tran.IPADDR, tran.TST}
	js, err := json.Marshal(transcationForVerify)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	decodedSignature, err := hex.DecodeString(tran.SIGNATURE)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	publicKey, err := hex.DecodeString(tran.SENDER)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	if secp.VerifySignature(js, decodedSignature, publicKey) != 1 {
		//

		statusCode = StatusSignVerifyError
		return
	}

	ok = true

	return
}

func ParseVoteTransaction(rawData string) (VoteTransaction, error) {
	//

	var voteTransaction VoteTransaction
	err := json.Unmarshal([]byte(rawData), &voteTransaction)
	return voteTransaction, err
}

func VerifyVoteTransaction(tran VoteTransaction) (transactionTime int64, statusCode int, ok bool) {
	//

	transactionTime = 0
	ok = false
	statusCode = StatusOk

	if len(tran.SENDER) != 66 {
		//

		statusCode = StatusWrongAttr_SENDER
		return
	}

	if len(tran.RECEIVER) != 66 {
		//

		statusCode = StatusWrongAttr_SENDER
		return
	}

	_, err := strconv.ParseInt(tran.VOTES, 10, 64)
	if err != nil {
		//

		statusCode = StatusWrongAttr_VOTES
		return
	}

	if len(tran.TST) != 10 {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	transactionTime, err = strconv.ParseInt(tran.TST, 10, 64)
	if err != nil {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	timestamp := time.Unix(transactionTime, 0)
	if int64(math.Abs(float64(time.Since(timestamp)/time.Second))) > 10 {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	if len(tran.SIGNATURE) != 130 {
		//

		statusCode = StatusWrongAttr_Signature
		return
	}

	transcationForVerify := VoteTransactionForVerify{tran.TT, tran.SENDER, tran.RECEIVER, tran.VOTES, tran.TST}
	js, err := json.Marshal(transcationForVerify)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	decodedSignature, err := hex.DecodeString(tran.SIGNATURE)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	publicKey, err := hex.DecodeString(tran.SENDER)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	if secp.VerifySignature(js, decodedSignature, publicKey) != 1 {
		//

		statusCode = StatusSignVerifyError
		return
	}

	ok = true

	return
}

func ParseUATransaction(rawData string) (UATransaction, error) {
	//

	var uaTransaction UATransaction
	err := json.Unmarshal([]byte(rawData), &uaTransaction)
	return uaTransaction, err
}

func VerifyUATransaction(tran UATransaction) (transactionTime int64, statusCode int, ok bool) {
	//

	transactionTime = 0
	ok = false
	statusCode = StatusOk

	if len(tran.SENDER) != 66 {
		//

		statusCode = StatusWrongAttr_SENDER
		return
	}

	if len(tran.TST) != 10 {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	transactionTime, err := strconv.ParseInt(tran.TST, 10, 64)
	if err != nil {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	timestamp := time.Unix(transactionTime, 0)
	if int64(math.Abs(float64(time.Since(timestamp)/time.Second))) > 10 {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	if len(tran.SIGNATURE) != 130 {
		//

		statusCode = StatusWrongAttr_Signature
		return
	}

	transcationForVerify := UATransactionForVerify{tran.TT, tran.SENDER, tran.TST}
	js, err := json.Marshal(transcationForVerify)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	decodedSignature, err := hex.DecodeString(tran.SIGNATURE)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	publicKey, err := hex.DecodeString(tran.SENDER)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	if secp.VerifySignature(js, decodedSignature, publicKey) != 1 {
		//

		statusCode = StatusSignVerifyError
		return
	}

	ok = true

	return
}

func ParseUVTransaction(rawData string) (UVTransaction, error) {
	//

	var uvTransaction UVTransaction
	err := json.Unmarshal([]byte(rawData), &uvTransaction)
	return uvTransaction, err
}

func VerifyUVTransaction(tran UVTransaction) (transactionTime int64, statusCode int, ok bool) {
	//

	transactionTime = 0
	ok = false
	statusCode = StatusOk

	if len(tran.SENDER) != 66 {
		//

		statusCode = StatusWrongAttr_SENDER
		return
	}

	if len(tran.RECEIVER) != 66 {
		//

		statusCode = StatusWrongAttr_SENDER
		return
	}

	if len(tran.TST) != 10 {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	transactionTime, err := strconv.ParseInt(tran.TST, 10, 64)
	if err != nil {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	timestamp := time.Unix(transactionTime, 0)
	if int64(math.Abs(float64(time.Since(timestamp)/time.Second))) > 10 {
		//

		statusCode = StatusWrongAttr_TST
		return
	}

	if len(tran.SIGNATURE) != 130 {
		//

		statusCode = StatusWrongAttr_Signature
		return
	}

	transcationForVerify := UVTransactionForVerify{tran.TT, tran.SENDER, tran.RECEIVER, tran.TST}
	js, err := json.Marshal(transcationForVerify)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	decodedSignature, err := hex.DecodeString(tran.SIGNATURE)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	publicKey, err := hex.DecodeString(tran.SENDER)
	if err != nil {
		//

		statusCode = StatusWrongDataFormat
		return
	}

	if secp.VerifySignature(js, decodedSignature, publicKey) != 1 {
		//

		statusCode = StatusSignVerifyError
		return
	}

	ok = true

	return
}

func CreateHelloTransaction(publicKey string, privateKey []byte, ip string) (tran string, ok bool) {
	//

	tran = ""
	ok = false

	helloTransactionForVerify := HelloTransactionForVerify{"HL", publicKey, ip, strconv.FormatInt(time.Now().Unix(), 10)}
	js, err := json.Marshal(helloTransactionForVerify)
	if err != nil {
		//

		return
	}

	signature := secp.Sign(js, []byte(privateKey))
	if len(signature) == 0 {
		//

		return
	}

	helloTransaction := HelloTransaction{"HL", publicKey, ip, helloTransactionForVerify.TST, string(signature)}
	js, err = json.Marshal(helloTransaction)
	if err != nil {
		//

		return
	}

	tran = string(js)
	ok = true

	return
}
