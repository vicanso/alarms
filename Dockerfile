FROM golang:1.14-alpine as builder

COPY . /alarms

RUN apk update \
  && apk add git make \
  && go get -u github.com/gobuffalo/packr/v2/packr2 \
  && cd /alarms \
  && make build

FROM alpine 

EXPOSE 7001

# tzdata 安装所有时区配置或可根据需要只添加所需时区

RUN addgroup -g 1000 go \
  && adduser -u 1000 -G go -s /bin/sh -D go \
  && apk add --no-cache ca-certificates tzdata

COPY --from=builder /alarms/alarms /usr/local/bin/alarms
COPY --from=builder /alarms/entrypoint.sh /home/go/entrypoint.sh

USER go

WORKDIR /home/go

HEALTHCHECK --timeout=10s CMD [ "wget", "http://127.0.0.1:7001/ping", "-q", "-O", "-"]

CMD ["alarms"]

ENTRYPOINT ["entrypoint.sh"]