# vim: filetype=dockerfile
FROM golang:1.14 as builder
ARG PKG
WORKDIR /go/src/${PKG}
ADD . .
RUN make build

FROM debian:buster-slim
RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /root/
COPY --from=builder /tmp/backend /usr/local/bin/backend
ENTRYPOINT ["backend"]
