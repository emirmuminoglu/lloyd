package red

type JSONInterface interface {
	MarshalJSON() ([]byte, error)
}

func (ctx *Ctx) JSONBlobResponse(body []byte, statusCode ...int) {
	ctx.Response.Header.SetContentType("application/json")

	if len(statusCode) > 0 {
		ctx.Response.Header.SetStatusCode(statusCode[0])
	}

	ctx.Response.SetBody(body)

	return
}

func (ctx *Ctx) JSONInterfaceResponse(data JSONInterface, statusCode ...int) {
	body, err := data.MarshalJSON()
	if err != nil {
		ctx.error = true
		return
	}

	ctx.JSONBlobResponse(body, statusCode...)

	return
}
