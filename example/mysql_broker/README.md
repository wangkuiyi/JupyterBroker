# How to Run This Example


## Build the Docker Image

Run the following command to build a Docker image that contains MySQL server, the `mysql` client command line tool, Go compiler, and our `mysql_broker` example:

```bash
docker build -t jupyterbroker:mysql .
```


## Run the Docker Image

Run the following command to start MySQL server and the `mysql_broker` server in a Docker container:

```bash
docker run --rm -d    \
    -p 3030:3030    \
	-e MYSQL_ROOT_PASSWORD=root   \
	-e MYSQL_ROOT_HOST='%'   \
	--name jm \
	jupyterbroker:mysql
```


## Try the Broker Server

Open a Web browser and direct to http://192.168.1.16:3030/mysql?sql=show%20databases%3B.  If you see a message like "MySQL server not yet started", wait few seconds and refresh.  You should see something like:

```
mysql: [Warning] Using a password on the command line interface can be insecure.
Database
information_schema
mysql
performance_schema
sys
```
