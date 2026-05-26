import { defineStore } from "pinia";
import { ref, computed } from "vue";
import { login, getUserInfo } from "@/api/auth";
import { ElMessage } from "element-plus";
import router from "@/router";

export const useUserStore = defineStore("user", () => {
  const token = ref<string>(localStorage.getItem("token") || "");
  const userInfo = ref<any>(JSON.parse(localStorage.getItem("userInfo") || "null"));

  const isLoggedIn = computed(() => !!token.value);
  const username = computed(() => userInfo.value?.username || "");
  const role = computed(() => userInfo.value?.role || "guest");

  // 登录
  const handleLogin = async (username: string, password: string) => {
    try {
      const res = await login({ username, password });
      token.value = res.token;
      userInfo.value = res.user;
      localStorage.setItem("token", res.token);
      localStorage.setItem("userInfo", JSON.stringify(res.user));
      ElMessage.success("登录成功！");
      router.push("/dashboard");
    } catch (err: any) {
      ElMessage.error(err.message || "登录失败");
    }
  };

  // 获取用户信息
  const fetchUserInfo = async () => {
    if (!token.value) return;
    try {
      const res = await getUserInfo();
      userInfo.value = res;
      localStorage.setItem("userInfo", JSON.stringify(res));
    } catch {
      logut();
    }
  };

  // 退出登录
  const logout = () => {
    token.value = "";
    userInfo.value = null;
    localStorage.removeItem("token");
    localStorage.removeItem("userInfo");
    router.push("/login");
  };

  return {
    token,
    userInfo,
    isLoggedIn,
    username,
    role,
    handleLogin,
    fetchUserInfo,
    logout,
  };
});
