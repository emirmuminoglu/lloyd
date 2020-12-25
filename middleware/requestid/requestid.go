package requestid

import (
	"github.com/emirmuminoglu/red"
	"github.com/google/uuid"
)

func NewV4() red.RequestHandler {
	return func(c *red.Ctx) {
		c.Request.Header.Set(red.XRequestIDHeader, uuid.New().String())
		c.Next()
	}
}

func NewV1() red.RequestHandler {
	return func(c *red.Ctx) {
		id, _ := uuid.NewUUID()
		c.Request.Header.Set(red.XRequestIDHeader, id.String())
		c.Next()
	}
}
