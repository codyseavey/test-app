FROM golang:alpine as builder
ENV GOPATH /go
ENV GOFLAGS '-mod=vendor'
ENV GO111MODULE on
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
RUN adduser -D -g '' testApp
WORKDIR ${GOPATH}/src/github.com/codyseavey/test-app
COPY . .
RUN GOOS=linux GOARCH=386 go build -ldflags="-w -s" -o /go/bin/testApp

FROM alpine
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/testApp .
USER testApp
EXPOSE 8080
ENTRYPOINT ["./testApp"]