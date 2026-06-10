FROM golang:alpine as builder

WORKDIR /go/src/kk

COPY . .
ENV GO111MODULE=on

RUN  go build -o ./server ./server


FROM golang:alpine
ENV TZ=Asia/Shanghai
WORKDIR /kk
RUN apk add --no-cache tzdata bash
COPY --from=builder /go/src/kk/server .
COPY --from=builder /go/src/kk/deploy/docker-compose/config/docker.config.yaml .

EXPOSE 8888
ENTRYPOINT exec ./server -c docker.config.yaml

CMD ["tail", "-f", "/dev/null"]