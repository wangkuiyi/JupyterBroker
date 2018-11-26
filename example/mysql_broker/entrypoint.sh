#!/bin/bash

$GOPATH/bin/mysql_broker -h 127.0.0.1 -p 3306 &
/entrypoint.sh mysqld

