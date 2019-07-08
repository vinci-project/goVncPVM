package restServer

import (
	"goVncPVM/helpers"
	"log"
	"strconv"

	"github.com/valyala/fasthttp"
)

func verifyPostTransaction(ctx *fasthttp.RequestCtx) (statusCode int,
	transactionForDB []byte,
	transactionTime int64) {
	//

	transactionForDB = ctx.PostBody()
	transactionForDBString := string(transactionForDB)
	log.Println(transactionForDBString)
	tranType := helpers.GetRawTransactionType(transactionForDBString)
	switch tranType {
	//

	case "ST":
		simpleTran, err := helpers.ParseSimpleTransaction(transactionForDBString)
		if err != nil {
			//

			return helpers.StatusWrongDataFormat, transactionForDB, transactionTime
		}

		transactionTime, statusCode, ok := helpers.VerifySimpleTransaction(simpleTran)
		if !ok {
			//

			return statusCode, transactionForDB, transactionTime
		}

		return helpers.StatusOk, transactionForDB, transactionTime

	return helpers.StatusUnknownTranType, transactionForDB, transactionTime
}

func verifyTranStatusRequest(ctx *fasthttp.RequestCtx) (statusCode int, key string) {
	//

	args := ctx.QueryArgs()
	for errNum, v := range helpers.RequestTranStatusFields {
		//

		if !args.Has(v) {
			//

			return errNum, key
		}
	}

	key = string(args.Peek("KEY"))

	if len(key) != 64 {
		//

		return helpers.StatusWrongAttr_KEY, key
	}

	return helpers.StatusOk, key
}

func verifyGetBlockRequest(ctx *fasthttp.RequestCtx) (statusCode int, bheight int64) {
	//

	args := ctx.QueryArgs()
	for errNum, v := range helpers.RequestGetBlockFields {
		//

		if !args.Has(v) {
			//

			return errNum, bheight
		}
	}

	heightString := string(args.Peek("BHEIGHT"))
	bheight, err := strconv.ParseInt(heightString, 10, 64)
	if err != nil {
		//

		return helpers.StatusWrongAttr_BHEIGHT, bheight
	}

	return helpers.StatusOk, bheight
}
