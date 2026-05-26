import request from "@/api/request";

// 获取监控统计
export function getMonitorStats() {
  return request.get("/monitor/stats");
}

// 获取流水线统计
export function getPipelineStats() {
  return request.get("/monitor/pipelines");
}

// 获取部署统计
export function getDeployStats() {
  return request.get("/monitor/deployments");
}

// 获取构建列表
export function getBuildList(params?: { page?: number; pageSize?: number }) {
  return request.get("/monitor/builds", { params });
}
