package mocksvc

import (
	"context"
	"encoding/json"
	"github.com/phayes/freeport"
	"net"
	"net/http"
	"strconv"
	"time"
)

type HandlerFunc func(req *http.Request) (int, interface{})

type MockService struct {
	handlerFunc HandlerFunc
	port        int
	server      *http.Server
}

func Start(handlerFunc HandlerFunc) *MockService {
	port, err := freeport.GetFreePort()
	if err != nil {
		panic("failed to get a free port for the mock service")
	}

	addr := ":" + strconv.Itoa(port)
	s := &MockService{
		handlerFunc: handlerFunc,
		port: port,
	}
	s.server = &http.Server{
		Addr:    addr,
		Handler: s,
	}

	go s.server.ListenAndServe()

	retries := 0
	var c net.Conn
	for {
		time.Sleep(time.Millisecond * 200)

		c, err = net.Dial("tcp4", s.server.Addr)
		if err != nil {
			if retries > 3 {
				panic("timed out waiting for the mock service to start")
			}

			retries++
		} else {
			_ = c.Close()
			break
		}
	}

	return s
}

func (s *MockService) Close() error {
	return s.server.Shutdown(context.Background())
}

func (s *MockService) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	status, res := s.handlerFunc(request)
	j, err := json.Marshal(res)
	if err != nil {
		panic("failed to marshal json: " + err.Error())
	}
	writer.WriteHeader(status)
	if _, err := writer.Write(j); err != nil {
		panic("failed to write the response: " + err.Error())
	}
}

func (s *MockService) BaseURL() string {
	return "http://localhost:" + strconv.Itoa(s.port)
}
