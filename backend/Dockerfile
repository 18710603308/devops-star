FROM golang:1.22-alpine as builder

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

COPY go.mod ./

# Alpine 镜像默认没有 git，go mod download 需要 git
RUN apk add --no-cache git

RUN go mod download

# 生成 go.sum（如果宿主机没有提交 go.sum）
RUN go mod tidy

COPY . .
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
