FROM mysql/mysql-server:8.0

RUN yum install -y tar
RUN curl https://dl.google.com/go/go1.11.2.linux-amd64.tar.gz | tar -C /usr/local -xzf -

ENV GOROOT /usr/local/go
ENV PATH $GOROOT/bin:$PATH
ENV GOPATH $HOME/go

RUN yum install -y git
RUN go get github.com/wangkuiyi/JupyterBroker/example/mysql_broker

EXPOSE 3030

COPY entrypoint.sh /broker-entrypoint.sh
ENTRYPOINT "/broker-entrypoint.sh"
