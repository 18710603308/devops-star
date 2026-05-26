import request from "@/api/request";

// 监控 API 对象
export const monitorAPI = {
  // 获取总统计
  getStats() {
    return request.get("/monitor/stats");
  },

  // 获取流水线统计
  getPipelineStats() {
    return request.get("/monitor/pipelines");
  },

  // 获取部署统计
  getDeployStats() {
    return request.get("/monitor/deployments");
  },

  // 获取构建列表
  getBuildList(params?: { page?: number; pageSize?: number }) {
    return request.get("/monitor/builds", { params });
  },
};

// 保留原有函数导出（兼容）
export function getMonitorStats() {
  return request.get("/monitor/stats");
}

export function getPipelineStats() {
  return request.get("/monitor/pipelines");
}

export function getDeployStats() {
  return request.get("/monitor/deployments");
}

export function getBuildList(params?: { page?: number; pageSize?: number }) {
  return request.get("/monitor/builds", { params });
}
