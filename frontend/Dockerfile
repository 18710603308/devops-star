FROM node:20-alpine as builder

WORKDIR /app

# 使用淘宝 npm 镜像加速
RUN npm config set registry https://registry.npmmirror.com

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# ========== 生产阶段 ==========
FROM nginx:alpine

WORKDIR /usr/share/nginx/html

COPY --from=builder /app/dist ./
COPY docker/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
