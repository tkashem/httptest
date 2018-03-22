package httpmock

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

// Server is a mocked http service
type Server interface {
	Setup() Call
	Close()
}

// Call allows a test to specify what to expect in an HTTP call and what to reply with
type Call interface {
	Expect(method string, path string, header http.Header) Call
	ReplyWith(statusCode int, header http.Header, body string, preAction func()) Call
	Remove()
}

// New returns a new instance of Server
func New(isTLS bool) Server {
	rootHandler := &handler{
		calls: make(map[*call]bool),
	}

	// Start the server
	s := httptest.NewUnstartedServer(rootHandler)

	switch isTLS {
	case true:
		go s.StartTLS()
	case false:
		go s.Start()
	}

	<-time.After(1 * time.Second)
	fmt.Printf("Mock http server has started, URL: %s\n", s.URL)

	return &server{
		handler: rootHandler,
		server:  s,
	}
}

type server struct {
	handler *handler
	server  *httptest.Server
}

func (s *server) Close() {
	if s.server != nil {
		fmt.Println("Closing the HTTP server")
		s.server.Close()
	}
}

func (s *server) Setup() Call {
	c := &call{
		handler: s.handler,
		request: &request{},
		reply:   &reply{},
	}

	c.handler.RegisterCall(c)
	return c
}
