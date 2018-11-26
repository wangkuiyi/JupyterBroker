package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	jupyterbroker "github.com/wangkuiyi/JupyterBroker"
)

type SQLRunner struct {
	jupyterbroker.ProcessRunner
}

var (
	mysqlHost = flag.String("h", "", "MySQL server host")
	mysqlPort = flag.Int("p", 0, "MySQL server port")
)

func NewSQLRunner(sql string) *SQLRunner {
	cmd := fmt.Sprintf("echo '%s' | mysql -uroot -proot -h %s -P %d",
		sql, *mysqlHost, *mysqlPort)
	return &SQLRunner{
		jupyterbroker.ProcessRunner{
			Cmd:  "sh",
			Args: []string{"-c", cmd}}}
}

func ProcessRunnerHandler(rw http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	sql, _ := url.QueryUnescape(req.Form["sql"][0])
	log.Println("Executing ", sql)
	NewSQLRunner(sql).Run(rw) // MakeSSEHandler will guard panics.
}

func main() {
	flag.Parse()
	http.HandleFunc("/mysql", jupyterbroker.MakeSSEHandler(ProcessRunnerHandler))
	http.ListenAndServe(":3030", nil)
}
