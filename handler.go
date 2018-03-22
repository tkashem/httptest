package httpmock

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ryanuber/go-glob"
)

type handler struct {
	calls map[*call]bool
}

func (h *handler) find(path string) *call {
	for c := range h.calls {
		if glob.Glob(c.request.path, path) {
			fmt.Printf("MATCHED:%s %s\n", c.request.path, path)
			return c
		}
		fmt.Printf("NOMATCHED:%s %s\n", c.request.path, path)
	}

	return nil
}

func (h *handler) handleReply(reply *reply, w http.ResponseWriter) {
	if reply.preAction != nil {
		reply.preAction()
	}

	if reply.statusCode != http.StatusOK {
		w.WriteHeader(reply.statusCode)
		fmt.Printf("REPLY: STATUS=%d\n", reply.statusCode)
	} else {
		for k, v := range reply.header {
			for _, value := range v {
				w.Header().Add(k, value)
			}
		}
		fmt.Printf("REPLY: HEADER=%v\n", w.Header())
	}
	w.Write([]byte(reply.body))
	bodySize := len(reply.body)
	if bodySize > 50 {
		bodySize = 50
	}
	fmt.Printf("REPLY: BODY=%s\n", reply.body[:bodySize])
}

func defaultReply() *reply {
	defaultReply := reply{
		statusCode: http.StatusNotFound,
	}
	return &defaultReply
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)

	fmt.Printf("METHOD=%s PATH=%s BODY=%s\n", r.Method, r.URL.Path, string(body))

	call := h.find(r.URL.Path)
	if call == nil {
		h.handleReply(defaultReply(), w)
		return
	}

	h.handleReply(call.reply, w)
}

func (h *handler) UnregisterCall(c *call) {
	delete(h.calls, c)
}

func (h *handler) RegisterCall(c *call) {
	h.calls[c] = true
}
