package nr

import (
	"context"

	"github.com/emirmuminoglu/lloyd"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const userValueName = "nr-txn"

func Middleware(app *newrelic.Application) lloyd.RequestHandler {
	return func(c *lloyd.Ctx) {
		method := lloyd.B2S(c.Request.Header.Method())
		url := c.URL()

		txn := app.StartTransaction(method + " " + c.PathName())

		var transport newrelic.TransportType

		switch url.Scheme {
		case "http":
			transport = newrelic.TransportHTTP
		case "https":
			transport = newrelic.TransportHTTPS
		default:
			transport = newrelic.TransportOther
		}

		txn.SetWebRequest(newrelic.WebRequest{
			URL:       url,
			Host:      url.Host,
			Method:    method,
			Transport: transport,
		})

		c.SetUserValue(userValueName, txn)
		c.Defer(func() {
			txn.SetWebResponse(c.ResponseWriter())
			txn.End()
		})
		c.Next()
	}
}

func GetTxn(ctx context.Context) *newrelic.Transaction {
	v := ctx.Value(userValueName)
	if v == nil {
		return nil
	}

	return v.(*newrelic.Transaction)
}
