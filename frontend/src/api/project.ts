import request from "@/api/request";

// 获取项目列表
export function getProjects() {
  return request.get("/projects");
}

// 创建项目
export function createProject(data: {
  name: string;
  display_name?: string;
  description?: string;
  repo_url?: string;
  repo_type?: string;
}) {
  return request.post("/projects", data);
}

// 获取项目详情
export function getProject(id: string) {
  return request.get(`/projects/${id}`);
}

// 更新项目
export function updateProject(id: string, data: any) {
  return request.put(`/projects/${id}`, data);
}

// 删除项目
export function deleteProject(id: string) {
  return request.delete(`/projects/${id}`);
}

// 获取项目成员
export function getProjectMembers(id: string) {
  return request.get(`/projects/${id}/members`);
}

// 添加项目成员
export function addProjectMember(id: string, data: { user_id: number; role: string }) {
  return request.post(`/projects/${id}/members`, data);
}
