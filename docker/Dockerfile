FROM ubuntu:14.04

# Note: The official Debian and Ubuntu images automatically ``apt-get clean``
# after each ``apt-get``
RUN apt-get update && apt-get install -y golang git 

RUN mkdir -p /go && \
    mkdir -p /go/src/github.com && \
    mkdir -p /go/bin && \
    mkdir -p /go/pkg

ENV GOPATH=/go
WORKDIR /go

RUN cd /go/src/github.com && \
    git clone https://github.com/cryptobase/crypto-scraper.git && \
    cd crypto-scraper/app && \
    go get && \
    go build && \
    mv app /usr/local/bin/crypto-scraper

RUN mkdir -p /data

VOLUME ["/data"]

CMD ["/usr/local/bin/crypto-scraper", "-path=/data/"]
