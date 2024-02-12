# 使用Go官方镜像作为构建环境
FROM golang:1.21.6-alpine3.19 as builder
# 创建工作目录
WORKDIR /app
# 设置Go代理，便于下载依赖
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct
# 先将go.mod和go.sum文件复制过去，然后下载依赖
# 先将全部的源代码复制过去
COPY . .
#编译
RUN go mod download  \
    && go build -trimpath -ldflags '-s -w -buildid=' -o main .

# 创建一个小型的最终镜像
FROM alpine:3.19
# 设置工作目录
WORKDIR /app
# 从构建镜像中复制编译好的应用程序
COPY --from=builder /app/main /app/main
# 设置标签

LABEL AUTHOR="guanren"
LABEL LANGUAGE="golang"
LABEL COPYRIGHT="guanren"
# 从构建镜像中复制编译好的应用程序，安装必要的包，并设置时区
# 安装ca-certificates以支持SSL，并设置时区为上海
RUN apk update --no-cache \
    && apk add --no-cache ca-certificates tzdata \
    && rm -rf /var/cache/apk/* \
    && mkdir -p /app/config \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# 设置时区环境变量
ENV TZ=Asia/Shanghai
EXPOSE 2095
# 定义容器启动时运行的命令
CMD ["/app/main"]
