# 构建阶段
FROM golang:1.22-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download
RUN go mod tidy

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux go build -o main .

# 运行阶段
FROM alpine:latest

# 安装 SQLite 运行时依赖
RUN apk add --no-cache sqlite-libs

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

CMD ["./main"]