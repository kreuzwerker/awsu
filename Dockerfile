FROM golang:1.15.13

RUN apt-get update && apt-get install -y libpcsclite-dev
RUN go get github.com/gesellix/awsu

CMD awsu
