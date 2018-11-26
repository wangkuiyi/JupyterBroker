# JupyterBroker

[![Build Status](https://travis-ci.org/wangkuiyi/JupyterBroker.svg?branch=develop)](https://travis-ci.org/wangkuiyi/JupyterBroker)

## Motivations

Suppose that we want to run SQL statements in Jupyter Notebook [cells](https://jupyter-notebook.readthedocs.io/en/stable/examples/Notebook/Notebook%20Basics.html?highlight=code%20cell#Edit-mode).  We can write a Jupyter Notebook kernel for SQL, or make it easier by reusing the existing IPython kernel but enabling it to run SQL by adding a [magic command](https://ipython.org/ipython-doc/3/config/custommagics.html#defining-magics).  Users mark a cell by `%%sql` and write the SQL code like the following:

```sql
%%sql
SELECT * FROM mydb.mytable;
```

The IPython kernel would call a Python class pre-configured and associate to the `%%sql` mark.  This magic class should send the SQL code to a broker server and copy the response to its stdout.  Jupyter Notebook server captures and displays magic command outputs.  The broker server, which usually runs in another Docker container or computer than the one running the Jupyter Notebook server, calls ODBC API to run the SQL code and streams outputs back to the magic class via [server-sent events](https://en.wikipedia.org/wiki/Server-sent_events).

Such a broker server is extensible to support other languages, Bash, AWK, Perl, name a few, and other tasks.  So I wrote this framework to ease the programming of the server in Go.


## The Framework

To define a broker and run it as a server, we need the following steps:

1. Add a runner type with the `Run(io.Writer)` method, which does the work and writes outputs to the writer.
1. Define an `http.HandlerFunc` function, which creates an instance of the runner type by parsing parameters from the HTTP request.
1. Register the handler with a URL using Go's standard `http.HandleFunc` API, and run the HTTP server.


## Example

This project includes an example runner, the `ProcessRunner`, which forks a sub-process to run a command line defined by `Cmd`, `Args`, and `Envs`.

```go
type ProcessRunner struct {
    Cmd        string
    Args, Envs []string
}
```

To write a SQL runner broker server, we can define an `SQLRunner` by using `ProcessRunner` to run a shell process, which runs the `echo <sql> | mysql` script.

```go
type SQLRunner struct {
    ProcessRunner
}

func NewSQLRunner(sql string) *SQLRunner {
    return &SQLRunner{
        Cmd: "sh", 
        Args: []string{"-c",
                       fmt.Sprintf("echo %s | mysql", sql)}}
}
```


To instantiate a SQL runner using parameters parsed from an HTTP request, we need to define a handler:

```go
func ProcessRunnerHandler(rw http.ResponseWriter, req *http.Request) {
    req.ParseForm()
    sql, _ := url.QueryUnescape(req.Form["sql"][0])
    NewSQLRunner(sql).Run(rw)  // MakeSSEHandler will guard panics.
}
```

To start the broker server, we need a simple `main` function:

```go
func main() {
    http.HandleFunc("/mysql", jupyterbroker.MakeSSEHandler(ProcessRunnerHandler))
    http.ListenAndServe(":3030", nil)
}
```

To test the server, type `http://localhost:3030/mysql?sql=SELECT%20*%20FROM*mydb.mytable` in your Web browser.
