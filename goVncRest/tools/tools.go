package tools

import (
	"encoding/json"
	"goVncPVM/helpers"

	"github.com/valyala/fasthttp"
)

func MakeResponse(statusCode int,
	ctx *fasthttp.RequestCtx) {
	//

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
}

func MakeDataResponse(data string,
	statusCode int,
	ctx *fasthttp.RequestCtx) {
	//

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
	ctx.SetBody([]byte(data))
}

func MakeVersionResponse(version string,
	statusCode int,
	ctx *fasthttp.RequestCtx) {
	//

	response := helpers.VersionResponse{version}
	jsResponse, _ := json.Marshal(response)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
	ctx.SetBody(jsResponse)
}

func MakeBHeightResponse(bheight string,
	statusCode int,
	ctx *fasthttp.RequestCtx) {
	//

	response := helpers.BHeightResponse{bheight}
	jsResponse, _ := json.Marshal(response)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
	ctx.SetBody(jsResponse)
}
