import request from "@/api/request";

// 获取部署环境列表
export function getEnvironments() {
  return request.get("/deploy/environments");
}

// 创建部署环境
export function createEnvironment(data: {
  name: string;
  display_name?: string;
  project_id: number;
  deploy_type?: string;
}) {
  return request.post("/deploy/environments", data);
}

// 触发部署
export function triggerDeploy(data: {
  environment_id: number;
  pipeline_run_id?: string;
  image_tag?: string;
}) {
  return request.post("/deploy/deploy", data);
}

// 获取部署历史
export function getDeployHistory(params?: {
  page?: number;
  page_size?: number;
}) {
  return request.get("/deploy/history", { params });
}
