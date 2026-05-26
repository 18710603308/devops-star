import axios, { AxiosInstance, AxiosResponse } from "axios";
import { ElMessage } from "element-plus";
import router from "@/router";
import { useUserStore } from "@/stores/user";

const baseURL =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

const request: AxiosInstance = axios.create({
  baseURL: `${baseURL}/api/v1`,
  timeout: 15000,
  headers: {
    "Content-Type": "application/json",
  },
});

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore();
    if (userStore.token) {
      config.headers!.Authorization = `Bearer ${userStore.token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    const res = response.data;
    // 如果是直接返回数据的接口（如登录）
    if (response.status === 200 && res.token) {
      return res;
    }
    if (res.code && res.code !== 0 && res.code !== 200) {
      ElMessage.error(res.message || "请求失败");
      return Promise.reject(new Error(res.message || "Error"));
    }
    // 有 message 字段时展示给用户（如 Harbor 未连接等提示）
    if (res.message) {
      ElMessage.warning(res.message);
    }
    return res;
  },
  (error) => {
    if (error.response) {
      const { status } = error.response;
      if (status === 401) {
        const userStore = useUserStore();
        userStore.logout();
        router.push("/login");
        ElMessage.error("登录已过期，请重新登录");
      } else if (status === 403) {
        ElMessage.error("无权限访问");
      } else if (status === 500) {
        ElMessage.error("服务器内部错误");
      } else {
        ElMessage.error(error.message || "网络错误");
      }
    } else {
      ElMessage.error("网络连接失败，请检查网络");
    }
    return Promise.reject(error);
  }
);

export default request;
