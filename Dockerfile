# 构建阶段
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 复制依赖文件并下载（利用缓存）
COPY go.mod go.sum ./
RUN go mod download

# 复制全部源码
COPY . .

# 编译静态二进制（禁用 CGO）
RUN CGO_ENABLED=0 GOOS=linux go build -o zhihu ./cmd/main.go

# 运行阶段
FROM alpine:latest

# 安装 CA 证书（如果需要访问外部 HTTPS）
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/zhihu .

# 暴露应用端口（根据你的实际端口修改，假设为 8080）
EXPOSE 8080

# 运行
CMD ["./zhihu"]