package jupyterbroker

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

type Runner interface {
	Run(w io.Writer)
}

func setServerSentEventHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func MakeSSEHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				http.Error(w, fmt.Sprintf("%v", e),
					http.StatusInternalServerError)
			}
		}()
		setServerSentEventHeader(w)
		f(w, r)
	}
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
