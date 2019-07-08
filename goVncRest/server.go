package restServer

import (
	"goVncPVM/goVncRest/tools"
	"goVncPVM/helpers"
	"log"
	"net"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/valyala/fasthttp"
)

var redisDB *redis.Client
var tranChannel *chan string

func fastHTTPRawHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetConnectionClose()
	if string(ctx.Method()) == "GET" {
		//

		switch string(ctx.Path()) {

		case "/pvm/tranStatus":
			statusCode, key := verifyTranStatusRequest(ctx)
			if statusCode >= 600 {
				//

				tools.MakeResponse(statusCode, ctx)
				return

			} else {
				//

				zScore := redisDB.ZScore("COMPLETE TRANSACTIONS", key)
				if !helpers.IsRedisError(zScore) {
					//

					tools.MakeResponse(helpers.StatusOk, ctx)
					return
				}

				zScore = redisDB.ZScore("FAILED TRANSACTIONS", key)
				if !helpers.IsRedisError(zScore) {
					//

					tools.MakeResponse(helpers.StatusTranFailed, ctx)
					return
				}

				tools.MakeResponse(helpers.StatusTranNotFound, ctx)
				return
			}

		case "/blockchain/getBHeight":
			intCmd := redisDB.ZCard("VNCCHAIN")
			if helpers.IsRedisError(intCmd) {
				//

				tools.MakeResponse(helpers.StatusInternalServerError, ctx)
				return
			}

			tools.MakeBHeightResponse(strconv.FormatInt(intCmd.Val(), 10), helpers.StatusOk, ctx)

		case "/blockchain/getTran":
			statusCode, key := verifyTranStatusRequest(ctx)
			if statusCode >= 600 {
				//

				tools.MakeResponse(statusCode, ctx)
				return

			} else {
				//

				stringCmd := redisDB.Get("TRANSACTIONS:" + key)
				if helpers.IsRedisError(stringCmd) {
					//

					tools.MakeResponse(helpers.StatusTranNotFound, ctx)
					return
				}

				tools.MakeDataResponse(stringCmd.Val(), helpers.StatusOk, ctx)
				return
			}

		case "/blockchain/getBlock":
			statusCode, bheight := verifyGetBlockRequest(ctx)
			if statusCode >= 600 {
				//

				tools.MakeResponse(statusCode, ctx)
				return

			} else {
				//

				stringSliceCmd := redisDB.ZRange("VNCCHAIN", bheight-1, bheight)
				if helpers.IsRedisError(stringSliceCmd) {
					//

					tools.MakeResponse(helpers.StatusDataNotFound, ctx)
					return
				}

				if len(stringSliceCmd.Val()) != 1 {
					//

					tools.MakeResponse(helpers.StatusDataNotFound, ctx)
					return
				}

				tools.MakeDataResponse(stringSliceCmd.Val()[0], helpers.StatusOk, ctx)
				return
			}

		case "/blockchain/getVersion":
			stringCmd := redisDB.Get("VERSION")
			if helpers.IsRedisError(stringCmd) {
				//

				tools.MakeResponse(helpers.StatusDataNotFound, ctx)
				return
			}

			tools.MakeVersionResponse(stringCmd.Val(), helpers.StatusOk, ctx)
			return

		case "/blockchain/getNodes":
			stringCmd := redisDB.Get("NODES LIST")
			if helpers.IsRedisError(stringCmd) {
				//

				tools.MakeResponse(helpers.StatusDataNotFound, ctx)
				return
			}

			tools.MakeDataResponse(stringCmd.Val(), helpers.StatusOk, ctx)
			return

		default:
			//

			ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		}

		return

	} else if string(ctx.Method()) == "POST" {
		//

		switch string(ctx.Path()) {

		case "/pvm/transaction":
			statusCode, transactionForDb, transactionTime := verifyPostTransaction(ctx)
			if statusCode >= 600 {
				//

				tools.MakeResponse(statusCode, ctx)
				return

			} else {
				//

				errRedis := redisDB.ZAdd("RAW TRANSACTIONS", redis.Z{
					Score:  float64(transactionTime),
					Member: transactionForDb,
				})

				if helpers.IsRedisError(errRedis) {
					//

					tools.MakeResponse(helpers.StatusInternalServerError, ctx)
					return
				}

				*tranChannel <- string(transactionForDb)

				tools.MakeResponse(helpers.StatusOk, ctx)
				return
			}
		}
	}

	ctx.Error("Unsupported method", fasthttp.StatusMethodNotAllowed)
}

func Start(r *redis.Client, c *chan string, ip string) {
	//

	redisDB = r
	tranChannel = c
	log.Pintln("START SERVER")
	// listener, err := reuseport.Listen("tcp4", net.JoinHostPort(ip, "5000"))
	// if err != nil {
	// 	log.Fatalf("error in reuseport listener: %s", err)
	// }

	server := &fasthttp.Server{
		Handler:          fastHTTPRawHandler,
		DisableKeepalive: true,
	}

	panic(server.ListenAndServe(net.JoinHostPort(ip, "5000")))
}
