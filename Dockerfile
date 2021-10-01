FROM golang as go 

COPY . /backend
WORKDIR /backend

RUN go env -w CGO_ENABLED=0 && go env -w GOPROXY=https://goproxy.cn,direct && go build -o proxy_server


FROM alpine:3.14

COPY --from=go /backend/proxy_server /IPProxy/proxy_server
COPY --from=go /backend/conf/config-docker.yaml /IPProxy/config.yaml

WORKDIR /IPProxy
EXPOSE 8000

CMD [ "./proxy_server" ]
