# 使用Go官方镜像作为构建环境
FROM golang:1.21.6-alpine3.19 as builder
# 创建工作目录
WORKDIR /app
# 设置Go代理，便于下载依赖
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct
# 先将go.mod和go.sum文件复制过去，然后下载依赖
COPY go.mod go.sum ./
# 下载依赖项
RUN go mod download

# 将本地代码复制到容器中
COPY . .
#编译
RUN go build -trimpath -ldflags '-s -w -buildid=' -o main .

# 创建一个小型的最终镜像
FROM alpine:3.19

# 设置工作目录
WORKDIR /app

# 安装ca-certificates以支持SSL，并设置时区为上海
RUN apk update --no-cache \
    && apk add --no-cache ca-certificates tzdata \
    && rm -rf /var/cache/apk/* \
    && mkdir -p /app/config \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# 从构建镜像中复制编译好的应用程序
COPY --from=builder /app/main .

# 设置时区环境变量
ENV TZ=Asia/Shanghai

# 定义容器启动时运行的命令
CMD ["/app/main"]
