FROM golang:alpine AS builder

WORKDIR /go/src/app

COPY . .

RUN go get

RUN go build -ldflags="-w -s" -o /go/bin/app

FROM alpine

COPY --from=builder /go/bin/app /go/bin/app

CMD ["/go/bin/app"]
