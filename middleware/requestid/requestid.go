package requestid

import (
	"github.com/emirmuminoglu/lloyd"
	"github.com/google/uuid"
)

func NewV4() lloyd.RequestHandler {
	return func(c *lloyd.Ctx) {
		c.Request.Header.Set(lloyd.XRequestIDHeader, uuid.New().String())
		c.Next()
	}
}

func NewV1() lloyd.RequestHandler {
	return func(c *lloyd.Ctx) {
		id, _ := uuid.NewUUID()
		c.Request.Header.Set(lloyd.XRequestIDHeader, id.String())
		c.Next()
	}
}
