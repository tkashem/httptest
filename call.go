package httpmock

import (
	"net/http"
)

type call struct {
	handler *handler
	request *request
	reply   *reply
}

type request struct {
	method string
	path   string
	header http.Header
}

type reply struct {
	statusCode int
	header     http.Header
	body       string
	preAction  func()
}

func (c *call) Remove() {
	c.handler.UnregisterCall(c)
}

func (c *call) Expect(method string, path string, header http.Header) Call {
	c.request.method = method
	c.request.header = header
	c.request.path = path

	return c
}

func (c *call) ReplyWith(statusCode int, header http.Header, body string, preAction func()) Call {
	c.reply.statusCode = statusCode
	c.reply.header = header
	c.reply.body = body
	c.reply.preAction = preAction

	return c
}
