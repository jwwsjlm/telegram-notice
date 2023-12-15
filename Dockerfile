FROM golang:1-alpine as builder

RUN mkdir /app

ADD . /app/
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags '-s -w -buildid=' -o main .

FROM alpine:latest

WORKDIR /app
RUN mkdir -p /app/config
COPY --from=builder /app/main .

CMD ["/app/main"]
