FROM golang:1.23-alpine as builder

WORKDIR /app

# 使用国内 Go Proxy 加速
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

# 清除可能从宿主机继承的代理环境变量（宿主机代理在容器内不可达）
ENV HTTP_PROXY=
ENV HTTPS_PROXY=
ENV http_proxy=
ENV https_proxy=
ENV ALL_PROXY=
ENV all_proxy=

# 先复制所有源文件（包括 go.mod、go.sum 如果本地有）
COPY . .

# Alpine 镜像默认没有 git，go mod tidy/download 需要 git
RUN apk add --no-cache git

# go mod tidy 自动清理无效依赖并生成/更新 go.sum
# go mod download 预下载依赖加速后续构建
RUN go mod tidy && go mod download

# 构建二进制（CGO_ENABLED=0 禁用 CGO，生成纯静态二进制）
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# ========== 生产阶段 ==========
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata curl

ENV TZ=Asia/Shanghai

WORKDIR /root/

COPY --from=builder /app/main .
COPY config/ config/
COPY scripts/ scripts/

EXPOSE 8080

HEALTHCHECK --interval=15s --timeout=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./main"]
