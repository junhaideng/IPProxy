FROM golang

COPY . /backend
WORKDIR /backend

RUN go env -w GOPROXY=https://goproxy.cn,direct && go build -o proxy_server main.go 

EXPOSE 8000

CMD [ "./proxy_server" ]
