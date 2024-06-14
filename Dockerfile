# 使用Go官方镜像作为构建环境，指定具体版本号
FROM golang:1.22.4-alpine as builder

# 创建工作目录
WORKDIR /app

# 安装upx用于压缩二进制文件，apk add命令执行后清理缓存
RUN apk update && apk add --no-cache upx && rm -rf /var/cache/apk/*

# 设置Go模块支持和代理，便于下载依赖
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 仅复制go.mod和go.sum，然后下载依赖，以加快构建速度并减小层的大小
COPY go.mod .
COPY go.sum .
RUN go mod download

# 复制全部的源代码
COPY . .

# 编译二进制文件，禁用CGO提高跨平台兼容性
# 使用-gcflags和-ldflags优化二进制文件大小和安全性
# 使用upx压缩main二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    && go build -gcflags "all=-N -l" -trimpath \
    -ldflags "-s -w -buildid=" -o main \
    && upx -6 main

# 创建一个小型的最终镜像
FROM alpine

# 设置工作目录
WORKDIR /app

# 从构建镜像中复制编译好的应用程序
COPY --from=builder /app/main /app/main

# 设置时区为上海
ENV TZ=Asia/Shanghai

# 使用标签来标记镜像的元数据
LABEL AUTHOR="guanren"
LABEL LANGUAGE="golang"
LABEL COPYRIGHT="guanren"

# 安装必要的包，设置时区，并清理缓存
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk update --no-cache \
    && apk add --no-cache tzdata ca-certificates libc6-compat libgcc libstdc++ curl \
    && rm -rf /var/cache/apk/* \
    && mkdir -p /app/config \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

# 声明应用程序运行时监听的端口
EXPOSE 2095

# 定义容器启动时运行的命令
CMD ["./main"]