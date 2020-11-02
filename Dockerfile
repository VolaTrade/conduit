FROM golang:1.12-alpine3.10 AS base
ARG GITHUB_TOKEN
RUN apk add bash ca-certificates git gcc g++ libc-dev make git
WORKDIR /go/src/github.com/volatrade/tickers
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
COPY config.env .
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/volatrade/".insteadOf "https://github.com/volatrade/" && go mod download


FROM base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -o bin/tickers


FROM alpine:3.10
RUN apk add ca-certificates
COPY --from=builder /go/src/github.com/volatrade/tickers/bin/tickers /bin/tickers
COPY --from=builder /go/src/github.com/volatrade/tickers/config.env .

CMD ["/bin/tickers"]
