FROM golang:1.17

MAINTAINER fsr

ENV GO111MODULE=on \
    CGO_ENABLE=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct"

WORKDIR /hm

COPY . .

# docker build时
RUN go get -u github.com/cosmtrek/air \
    && go mod download
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

EXPOSE 8080

# docker run时
#ENTRYPOINT ["air"]

