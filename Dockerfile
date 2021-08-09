FROM golang:1.16-alpine3.14 as build
RUN apk add --update make
WORKDIR /go/src/app
COPY . .
RUN make build
RUN make install
ENTRYPOINT ["tseep"]