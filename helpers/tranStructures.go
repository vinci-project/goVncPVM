package helpers

const StatusOk = 200 // OK

const (
	StatusNotEnoughFunds              = 600 + iota // NOT ENOUTH FUNDS
	StatusNotEnoughFundsForCommission              // 601 NOT ENOUGH FUNDS FOR COMMISSION
	StatusTranNotFound                             // 602 TRANSACTION NOT FOUND
	StatusInternalServerError                      // 603 SOMETHING BAD HAPPEND
	StatusUnknownTranType                          // 604 UNKNOWN TRANSACTION TYPE
	StatusAttrNotFound_TT                          // 605 CAN NOT FIND ATTRIBUTE - TT
	StatusAttrNotFound_VERSION                     // 606 CAN NOT FIND ATTRIBUT - VERSION
	StatusAttrNotFound_SENDER                      // 607 CAN NOT FIND ATTRIBUTE - SENDER
	StatusAttrNotFound_RECEIVER                    // 608 CAN NOT FIND ATTRIBUTE - RECEIVER
	StatusAttrNotFound_TTOKEN                      // 609 CAN NOT FIND ATTRIBUTE - TTOKEN
	StatusAttrNotFound_CTOKEN                      // 610 CAN NOT FIND ATTRIBUTE - CTOKEN
	StatusAttrNotFound_TST                         // 611 CAN NOT FIND ATTRIBUTE - TST
	StatusAttrNotFound_SIGNATURE                   // 612 CAN NOT FIND ATTRIBUTE - SIGNATURE
	StatusAttrNotFound_KEY                         // 613 CAN NOT FIND ATTRIBUTE - KEY
	StatusAttrNotFound_SOURCE                      // 614 CAN NOT FIND ATTRIBUTE - SOURCE
	StatusAttrNotFound_PAIR                        // 615 CAN NOT FIND ATTRIBUTE - PAIR
	StatusAttrNotFound_TICKER                      // 616 CAN NOT FIND ATTRIBUTE - TICKER
	StatusWrongAttr_KEY                            // 617 WRONG ATTRIBUTE - KEY
	StatusWrongAttr_TT                             // 618 WRONG ATTRIBUTE - TT
	StatusWrongAttr_VERSION                        // 619 WRONG ATTRIBUTE - VERSION
	StatusWrongAttr_SENDER                         // 620 WRONG ATTRIBUTE - SENDER
	StatusWrongAttr_RECEIVER                       // 621 WRONG ATTRIBUTE - RECEIVER
	StatusWrongAttr_TST                            // 622 WRONG ATTRIBUTE - TST
	StatusWrongAttr_CTOKEN                         // 623 WRONG ATTRIBUTE - CTOKEN
	StatusWrongAttr_Signature                      // 624 WRONG ATTRIBUTE - SIGNATURE
	StatusSignVerifyError                          // 625 CAN'T VERIFY SIGNATURE
	StatusDontSendYourself                         // 626 YOU CAN'T SEND YOURSELF
	StatusTranFailed                               // 627 TRANSACTION FAILED
	StatusDataNotFound                             // 628 DATA NOT FOUND
	StatusWrongDataFormat                          // 629 WRONG DATA FORMAT
	StatusWrongAttr_TTOKEN                         // 630 WRONG ATTRIBUTE - TTOKEN
	StatusWrongAttr_BHEIGHT                        // 631 WRONG ATTRIBUTE - BHEIGHT
	StatusWrongAttr_IPADDR                         // 632 WRONG ATTRIBUTE - IPADDR
	StatusWrongAttr_VOTES                          // 633 WRONG ATTRIBUTE - VOTES
	StatusAttrNotFound_BHEIGHT                     // 634 CAN NOT FIND ATTRIBUTE - BHEIGHT
	StatusAttrNotFound_IPADDR                      // 635 CAN NOT FIND ATTRIBUTE - IPADDR
	StatusAttrNotFound_ADDRESS                     // 636 CAN NOT FIND ATTRIBUTE - ADDRESS
	StatusWrongAttr_ADDRESS                        // 637 WRONG ATTRIBUTE - ADDRESS
)

var MandatoryTransactionFields = map[int]string{StatusWrongAttr_TT: "TT",
	//	StatusAttrNotFound_VERSION:   "VERSION",
	StatusAttrNotFound_SENDER:    "SENDER",
	StatusAttrNotFound_TST:       "TST",
	StatusAttrNotFound_SIGNATURE: "SIGNATURE"}

var SimpleStructureFields = map[int]string{StatusWrongAttr_TT: "TT",
	//	StatusAttrNotFound_VERSION:   "VERSION",
	StatusAttrNotFound_SENDER:    "SENDER",
	StatusAttrNotFound_RECEIVER:  "RECEIVER",
	StatusAttrNotFound_TTOKEN:    "TTOKEN",
	StatusAttrNotFound_CTOKEN:    "CTOKEN",
	StatusAttrNotFound_TST:       "TST",
	StatusAttrNotFound_SIGNATURE: "SIGNATURE"}

var ApplicantStructureFields = map[int]string{StatusWrongAttr_TT: "TT",
	//	StatusAttrNotFound_VERSION:   "VERSION",
	StatusAttrNotFound_SENDER:    "SENDER",
	StatusAttrNotFound_IPADDR:    "IPADDR",
	StatusAttrNotFound_TST:       "TST",
	StatusAttrNotFound_SIGNATURE: "SIGNATURE"}

var RequestBalanceFields = map[int]string{StatusAttrNotFound_TTOKEN: "TTOKEN",
	StatusAttrNotFound_SENDER: "SENDER"}

var RequestStakeFields = map[int]string{StatusAttrNotFound_SENDER: "SENDER"}

var RequestTranStatusFields = map[int]string{StatusAttrNotFound_KEY: "KEY"}

var RequestVSFields = map[int]string{StatusAttrNotFound_ADDRESS: "ADDRESS"}

var RequestGetBlockFields = map[int]string{StatusAttrNotFound_BHEIGHT: "BHEIGHT"}

type BalanceResponse struct {
	BALANCE map[string]string
	STAKE map[string]string
}

type StakeResponse struct {
	STAKE map[string]string
}

type VersionResponse struct {
	VERSION string
}

type BHeightResponse struct {
	BHEIGHT string
}

type ASResponse struct {
	APPLICANTS []string
}

type VSResponse struct {
	VOTES string
}

type AVSResponse struct {
	VOTES map[string]string
}

type SimpleTransaction struct {
	TT        string
	SENDER    string
	RECEIVER  string
	TTOKEN    string
	CTOKEN    string
	TST       string
	SIGNATURE string
}

type SimpleTransactionForVerify struct {
	TT       string
	SENDER   string
	RECEIVER string
	TTOKEN   string
	CTOKEN   string
	TST      string
}

type ApplicantTransaction struct {
	TT        string
	SENDER    string
	IPADDR    string
	TST       string
	SIGNATURE string
}

type ApplicantTransactionForVerify struct {
	TT     string
	SENDER string
	IPADDR string
	TST    string
}

type VoteTransaction struct {
	TT        string
	SENDER    string
	RECEIVER  string
	VOTES     string
	TST       string
	SIGNATURE string
}

type VoteTransactionForVerify struct {
	TT       string
	SENDER   string
	RECEIVER string
	VOTES    string
	TST      string
}

type UATransaction struct {
	TT        string
	SENDER    string
	TST       string
	SIGNATURE string
}

type UATransactionForVerify struct {
	TT     string
	SENDER string
	TST    string
}

type UVTransaction struct {
	TT        string
	SENDER    string
	RECEIVER  string
	TST       string
	SIGNATURE string
}

type UVTransactionForVerify struct {
	TT       string
	SENDER   string
	RECEIVER string
	TST      string
}

type TranType struct {
	TT string
}

type HelloTransaction struct {
	TT        string
	SENDER    string
	ADDRESS   string
	TST       string
	SIGNATURE string
}

type HelloTransactionForVerify struct {
	TT      string
	SENDER  string
	ADDRESS string
	TST     string
}

type StatisticsResponse struct {
	TCOUNT  string
	LTPS    string
	BHEIGHT string
	LTPB    string
	TPD     string
	VMAP    []string
	UPD     string
}
