package jupyterbroker

import (
	"log"
	"net"
	"net/http"
)

type ProcessBroker struct {
}

func setServerSentEventHeader(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
}

func (br *ProcessBroker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	setServerSentEventHeader(rw)

	req.ParseForm()
	pr := NewProcessRunner(req.Form["cmd"][0], req.Form["args"], req.Form["envs"])
	pr.Run(rw)
}

func Start(addr string) (string, error) {
	lst, e := net.Listen("tcp", addr)
	if e != nil {
		return "", e
	}
	go func() {
		log.Fatal("HTTP server error: ", http.Serve(lst, nil))
	}()
	return lst.Addr().String(), nil
}
