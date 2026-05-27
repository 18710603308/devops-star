import request from "@/api/request";

// 获取镜像仓库列表
export function getImages() {
  return request.get("/api/v1/artifacts/images");
}

// 获取镜像标签列表
export function getTags(repoName: string) {
  return request.get(`/api/v1/artifacts/images/${repoName}/tags`);
}

// 删除镜像（按标签）
export function deleteImage(repoName: string, tag: string) {
  return request.delete(`/api/v1/artifacts/images/${repoName}/${tag}`);
}

// 扫描镜像漏洞
export function scanImage(repoName: string, tag: string) {
  return request.post(`/api/v1/artifacts/images/${repoName}/${tag}/scan`);
}

// 获取镜像扫描报告
export function getScanReport(repoName: string, tag: string) {
  return request.get(`/api/v1/artifacts/images/${repoName}/${tag}/scan-report`);
}

// 获取 Harbor 连接状态
export function getStatus() {
  return request.get("/api/v1/artifacts/status");
}
