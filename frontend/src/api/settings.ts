import request from "@/api/request";

// 获取镜像源配置
export function getMirrorConfig() {
  return request.get("/settings/mirrors");
}

// 更新镜像源配置
export function updateMirrorConfig(data: Record<string, any>) {
  return request.put("/settings/mirrors", data);
}
