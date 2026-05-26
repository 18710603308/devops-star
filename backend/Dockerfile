FROM golang:1.22-alpine as builder

WORKDIR /app

# 使用国内 Go Proxy 加速
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

COPY go.mod go.sum ./
RUN go mod download

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
