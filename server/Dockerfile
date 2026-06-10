FROM golang:alpine as builder

WORKDIR /go/src/mock-subscribe

COPY . .
ENV GO111MODULE=on

RUN  go build -o ./server ./server


FROM golang:alpine
ENV TZ=Asia/Shanghai
WORKDIR /mock-subscribe
RUN apk add --no-cache tzdata bash
COPY --from=builder /go/src/mock-subscribe/server .
COPY --from=builder /go/src/mock-subscribe/server/config.docker.yaml ./config.docker.yaml

EXPOSE 8888
ENTRYPOINT exec ./server -c config.docker.yaml

CMD ["tail", "-f", "/dev/null"]