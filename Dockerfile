FROM golang:1.16-alpine3.14 as build
RUN apk add --update make
WORKDIR /go/src/app
COPY . .
RUN make build
RUN make install

FROM alpine:3.14 as run
RUN apk add --update iptables
WORKDIR /app
COPY --from=build /usr/local/bin/tseep .
ENTRYPOINT ["/app/tseep"]