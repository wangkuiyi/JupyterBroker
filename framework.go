package jupyterbroker

import (
	"log"
	"net/http"
)

type Broker struct {
	Runner
}

func NewBroker(r Runner) *Broker {
	return &Broker{r}
}

func (br *Broker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	br.Run(rw)
}

func Run() {
	http.Handle("/process", NewBroker(NewProcessRunner(
		"sh", []string{"-c", "echo hello $foo"}, []string{"foo=bar"})))
	log.Fatal("HTTP server error: ", http.ListenAndServe("localhost:3000", nil))
}
