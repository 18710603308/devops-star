import request from "@/api/request";

// 获取流水线列表
export function getPipelines(params?: { project_id?: number }) {
  return request.get("/pipelines", { params });
}

// 创建流水线
export function createPipeline(data: {
  name: string;
  description?: string;
  project_id: number;
  config_yaml?: string;
}) {
  return request.post("/pipelines", data);
}

// 获取流水线详情
export function getPipeline(id: string) {
  return request.get(`/pipelines/${id}`);
}

// 更新流水线
export function updatePipeline(id: string, data: {
  name?: string;
  description?: string;
  config_yaml?: string;
}) {
  return request.put(`/pipelines/${id}`, data);
}

// 删除流水线
export function deletePipeline(id: string) {
  return request.delete(`/pipelines/${id}`);
}

// 触发流水线
export function triggerPipeline(id: string) {
  return request.post(`/pipelines/${id}/trigger`);
}

// 获取流水线日志
export function getPipelineLogs(id: string) {
  return request.get(`/pipelines/${id}/logs`);
}

// 获取环境列表
export function getEnvironments() {
  return request.get("/deploy/environments");
}

// 创建环境
export function createEnvironment(data: {
  name: string;
  display_name?: string;
  project_id: number;
  deploy_type?: string;
}) {
  return request.post("/deploy/environments", data);
}

// 部署
export function deploy(data: {
  environment_id: number;
  pipeline_run_id?: string;
  image_tag?: string;
}) {
  return request.post("/deploy/deploy", data);
}

// 获取部署历史
export function getDeployHistory(params?: { page?: number; page_size?: number }) {
  return request.get("/deploy/history", { params });
}
