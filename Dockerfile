FROM golang:alpine AS builder

WORKDIR /go/src/app

COPY . .

RUN go get -u ./...

RUN go build -C cmd/link-forge/ -ldflags="-w -s" -o /go/bin/app

FROM alpine

COPY --from=builder /go/bin/app /go/bin/app

CMD ["/go/bin/app"]
