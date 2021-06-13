# 打包依赖阶段使用golang作为基础镜像
FROM golang:1.16 as build

# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /go/release

COPY . .

# 指定OS等，并go build
#RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o app main.go
RUN GOOS=linux go build -ldflags='-s -w' -o app .

# 运行阶段指定 busybox 作为基础镜像
FROM busybox

COPY --from=build /go/release/app /
COPY --from=build /go/release/conf.ini /

# 指定运行时环境变量
ENV GIN_MODE=release \
    PORT=80

EXPOSE 80
VOLUME /log

#ENTRYPOINT /app