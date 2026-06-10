FROM golang:alpine AS builder

WORKDIR /go/src/mock-subscribe

COPY server/ ./
ENV GO111MODULE=on

RUN go build -o /tmp/server .

FROM golang:alpine
ENV TZ=Asia/Shanghai
WORKDIR /mock-subscribe
RUN apk add --no-cache tzdata bash
COPY --from=builder /tmp/server ./server
COPY --from=builder /go/src/mock-subscribe/server/config.docker.yaml ./config.docker.yaml

EXPOSE 8888
ENTRYPOINT ["./server", "-c", "config.docker.yaml"]
