FROM alpine:3.12 as alpine

RUN apk add -U --no-cache ca-certificates

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY ./bin/conduit /bin/conduit
COPY config.env .


CMD ["/bin/conduit"]
