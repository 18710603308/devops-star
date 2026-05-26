import request from "@/api/request";

// 用户登录
export function login(data: { username: string; password: string }) {
  return request.post("/auth/login", data);
}

// 用户注册
export function register(data: { username: string; password: string; email: string }) {
  return request.post("/auth/register", data);
}

// 获取当前用户信息
export function getUserInfo() {
  return request.get("/auth/me");
}
