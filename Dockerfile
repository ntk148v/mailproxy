FROM golang:1.12-alpine as builder

ENV GO111MODULE=on
ENV APPLOC=$GOPATH/src/mailproxy

RUN apk add --no-cache git

ADD . $APPLOC
WORKDIR $APPLOC
RUN go build -mod vendor -o /bin/mailproxy

FROM golang:1.12-alpine
LABEL maintainer="Kien Nguyen <kiennt2609@gmail.com>"
COPY --from=builder /bin/mailproxy /bin/mailproxy
RUN chmod +x /bin/mailproxy && \
    mkdir /etc/mailproxy
EXPOSE 9011
ENTRYPOINT ["/bin/mailproxy"]
CMD ["-conf", "/etc/mailproxy"]
